package arm64

import (
	"errors"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// ByElem — by-element multiplications (MLA/UMULL/SQDMULH/FMUL families):
// the result = a vector times a lane of another vector. The common
// space [29:24]=x01111 is shared with shift-by-immediate via bit 10
// (1 = shift, 0 = by-element):
//
//	name    = (U=bit29) × (size=bits23:22: 01/10=int, 00/11=fp) × (opc=bits15:12)
//	index   = the top (4-?) bits of the field {b11,b21,b20,b19} (weights 8:4:2:1),
//	          MSB-aligned by the lane count: .h→3 bits, .s→2, .d→1
//	Vm      = bits[19:16] for a .h source, [20:16] for .s/.d
//	fcmla   = rotation from opc: 0001→#0, 0011→#90, 0101→#180, 0111→#270
type ByElem struct {
	base

	name          string
	q, size, idx  uint32
	rot           uint32 // fcmla: 0/90/180/270
	rd, rn, rm    string
	rdN, rnN, rmN uint32
	long          bool // long forms (smlal/...): result one size wider, Q → "2"
}

var byElemIntNames = [2][16]string{
	// U = 0 (signed)
	{1: "smlal", 3: "sqdmlal", 5: "smlsl", 7: "sqdmlsl",
		8: "mul", 0xa: "smull", 0xb: "sqdmull", 0xc: "sqdmulh", 0xd: "sqrdmulh"},
	// U = 1 (unsigned / accumulate)
	{0: "mla", 1: "fcmla", 2: "umlal", 3: "fcmla", 4: "mls", 5: "fcmla",
		6: "umlsl", 7: "fcmla", 0xa: "umull", 0xd: "sqrdmlah", 0xf: "sqrdmlsh"},
}

var byElemFPNames = [2][16]string{
	{1: "fmla", 5: "fmls", 9: "fmul"},
	{9: "fmulx"},
}

var byElemLong = map[string]bool{
	"smlal": true, "sqdmlal": true, "smlsl": true, "sqdmlsl": true,
	"smull": true, "sqdmull": true, "umlal": true, "umlsl": true, "umull": true,
}

func decodeByElem(w uint32, addr uint64) Instr {
	u, size, opc, q := w>>29&1, w>>22&3, w>>12&0xf, w>>30&1
	var name string
	if size == 1 || size == 2 {
		name = byElemIntNames[u][opc]
	} else {
		name = byElemFPNames[u][opc]
	}

	long := byElemLong[name]
	// index: MSB-aligned bits of the field {b11,b21,b20,b19}
	lanes := 4 - size // .b→4 .h→3 .s→2 .d→1
	field := w>>11&1<<3 | w>>21&1<<2 | w>>20&1<<1 | w>>19&1
	idx := field >> (4 - lanes)
	var rmN uint32
	if size == 1 {
		rmN = w >> 16 & 0xf // .h source: Vm is 4-bit
	} else {
		rmN = w >> 16 & 0x1f
	}

	name2 := name
	if long && q == 1 {
		name2 += "2"
	}

	return ByElem{
		base: newBase(addr, w),
		name: name2,
		q:    q,
		size: size,
		idx:  idx,
		rot:  opc / 2 * 90 % 360, // fcmla: 0001→0, 0011→90, 0101→180, 0111→270
		rd:   vReg(w & 0x1f),
		rn:   vReg(w >> 5 & 0x1f),
		rm:   vReg(rmN),
		rdN:  w & 0x1f,
		rnN:  w >> 5 & 0x1f,
		rmN:  rmN,
		long: long,
	}
}

func (i ByElem) ObjDump(_ disasm.ViewCtx) string {
	var arr string
	switch {
	case i.size == 0:
		arr = decodeArrangement(i.q, 2) // fp32: .2s/.4s
	case i.size == 3:
		arr = "2d" // fp64
	case i.long:
		// long forms: the result's lanes come from the Q=1 row (.2s→.2d,
		// .4h→.4s; Q=1 itself gives the "2" form)
		arr = decodeArrangement(1, i.size+1)
	default:
		arr = decodeArrangement(i.q, i.size)
	}

	s := fmt.Sprintf("%s.%s %s, %s, %s[%d]", i.name, arr, i.rd, i.rn, i.rm, i.idx)
	if i.name == "fcmla" {
		s += fmt.Sprintf(", #%d", i.rot)
	}

	return s
}

func (i ByElem) Encode(w io.Writer, pc uint64) (int64, error) {
	name := i.name
	if i.long && i.q == 1 {
		name = name[:len(name)-1]
	}

	u := uint32(0)
	if i.name[0] == 'u' || i.name == "mla" || i.name == "mls" || i.name == "fcmla" ||
		i.name == "fmulx" {
		u = 1
	}

	var opc uint32
	switch {
	case i.name == "fcmla":
		opc = i.rot/90*2 + 1 // 0→0001, 90→0011, 180→0101, 270→0111
	case i.size == 1 || i.size == 2:
		for c, n := range byElemIntNames[u] {
			if n == name {
				opc = uint32(c)
				break
			}
		}
	default:
		for c, n := range byElemFPNames[u] {
			if n == name {
				opc = uint32(c)
				break
			}
		}
	}

	f := i.idx << i.size // MSB-aligned in the 4-bit field {b11,b21,b20,b19}
	return writeWord(w, i.q<<30|u<<29|0x0F<<24|i.size<<22|i.rmN<<16|opc<<12|
		(f>>3&1)<<11|(f>>2&1)<<21|(f>>1&1)<<20|(f&1)<<19|i.rnN<<5|i.rdN)
}

func (i ByElem) MarshalJSON() ([]byte, error) {
	return i.marshal(i.name, i.ObjDump(disasm.DefaultViewCtx()), "ASIMD",
		map[string]any{"Rd": i.rd, "Rn": i.rn, "Rm": i.rm, "index": i.idx})
}

func init() {
	// int: long × {"", "2"} × 3 arrangements; the rest × 2q × size
	for u := range 2 {
		for _, name := range byElemIntNames[u] {
			if name == "" || name == "fcmla" {
				continue
			}

			if byElemLong[name] {
				// Q=0 and Q=1 forms print the same result
				// arrangement; the keys differ by the "2" suffix
				for size := range uint32(3) {
					arr := decodeArrangement(1, size+1)
					armCtors[name+"."+arr] = newByElemArm(name, 0, size, true)
					armCtors[name+"2."+arr] = newByElemArm(name, 1, size, true)
				}
			} else {
				for q := range uint32(2) {
					for size := uint32(1); size < 3; size++ {
						k := name + "." + decodeArrangement(q, size)
						armCtors[k] = newByElemArm(name, q, size, false)
					}
				}
			}
		}
	}

	// fp: .2s/.4s (size 0), .2d (size 3)
	for u := range 2 {
		for _, name := range byElemFPNames[u] {
			for q := range uint32(2) {
				k := name + "." + decodeArrangement(q, 2)
				armCtors[k] = newByElemArm(name, q, 0, false)
			}

			armCtors[name+".2d"] = newByElemArm(name, 1, 3, false)
		}
	}

	// fcmla: 2 arrangement classes (rotation — the #rot operand, opc computed
	// at encoding time)
	for q := range uint32(2) {
		for size := uint32(1); size < 3; size++ {
			armCtors["fcmla."+decodeArrangement(q, size)] = newByElemArm("fcmla", q, size, false)
		}
	}
}

// newByElemArm — name.Arr vd, vn, vm[idx]{, #rot (fcmla)}.
func newByElemArm(name string, q, size uint32, long bool) func([]vOp) (Instr, error) {
	return func(ops []vOp) (Instr, error) {
		if len(ops) < 3 || ops[0].arr == "" {
			return nil, fmt.Errorf("%s: want vd.Arr, vn, vm[idx]", name)
		}

		rd, err := wantV(ops[0], name)
		if err != nil {
			return nil, err
		}

		rn, err := wantV(ops[1], name)
		if err != nil {
			return nil, err
		}

		rm, err := wantV(ops[2], name)
		if err != nil {
			return nil, err
		}

		idx := ops[2].num
		if idx < 0 {
			return nil, fmt.Errorf("%s: bad index", name)
		}

		rdN, rnN, rmN, err := regNums3(rd, rn, rm)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}

		name2 := name
		if long && q == 1 {
			name2 += "2"
		}

		var rot uint32
		if name == "fcmla" {
			if len(ops) < 4 || ops[3].kind != armOpImm {
				return nil, errors.New("fcmla: rotation expected (#0/#90/#180/#270)")
			}

			r := ops[3].num
			if r%90 != 0 || r < 0 || r > 270 {
				return nil, errors.New("fcmla: bad rotation")
			}

			rot = uint32(r)
		}

		return ByElem{
			name: name2,
			q:    q,
			size: size,
			idx:  uint32(idx),
			rot:  rot,
			rd:   rd,
			rn:   rn,
			rm:   rm,
			rdN:  rdN,
			rnN:  rnN,
			rmN:  rmN,
			long: long,
		}, nil
	}
}
