package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// SimdCopy - the Advanced SIMD copy family: dup.Arr vd, wn (DUP general) |
// mov.sz vd[idx], wn (INS general) | smov/umov wd, vn[idx]. imm5 encodes
// size+index: size = the number of the lowest set bit, index = imm5 >>
// (size+1); the operation is in bits [13:12] (0=dup, 1=ins, 2=smov, 3=umov),
// bits [15:14]=00, bits [11:10]=11.
type SimdCopy struct {
	base

	op      string // dup | ins | smov | umov
	size    uint32 // 0=b 1=h 2=s 3=d
	idx     uint32
	q       uint32
	vd, gpr string // vd - the vector one, gpr - the x/w on the other side
	enc     uint32
	vdNum   uint32
	gprNum  uint32
	isDest  bool // gpr on the destination side (smov/umov)
}

func decodeSimdCopy(w uint32, addr uint64) Instr {
	imm5 := w >> 16 & 0x1f
	// imm5 must be one-hot (1/2/4/8 - the b/h/s/d size); the schema mask
	// does not express this, unencoded values go to the .word fallback.
	if imm5 == 0 || imm5 > 8 || imm5&(imm5-1) != 0 {
		return decodeUnknown(w, addr)
	}

	size := uint32(bitsCtz(imm5))
	idx := imm5 >> (size + 1)
	q := w >> 30 & 1
	var op string
	isDest := false
	switch w >> 12 & 3 { // bits [13:12]: 0=dup 1=ins 2=smov 3=umov
	case 0:
		op = "dup"
	case 1:
		op = "ins"
	case 2:
		op, isDest = "smov", true
	default:
		op, isDest = "umov", true
	}

	vdNum, gprNum := w&0x1f, w>>5&0x1f
	if isDest {
		// smov/umov: the GPR destination is in the Rd slot, the vector source in Rn
		vdNum, gprNum = w>>5&0x1f, w&0x1f
	}

	gpr64 := size == 3 || (op == "smov" && q == 1)
	vd := vReg(vdNum)
	gpr := armRegName(gprNum, gpr64)
	return SimdCopy{
		base:   newBase(addr, w),
		op:     op,
		size:   size,
		idx:    idx,
		q:      q,
		vd:     vd,
		gpr:    gpr,
		enc:    0x0E000000,
		vdNum:  vdNum,
		gprNum: gprNum,
		isDest: isDest,
	}
}

// bitsCtz - the number of the lowest set bit.
func bitsCtz(v uint32) int {
	n := 0
	for v&1 == 0 && v != 0 {
		n++
		v >>= 1
	}

	return n
}

func (i SimdCopy) ObjDump(_ disasm.ViewCtx) string {
	sz := [...]string{"b", "h", "s", "d"}[i.size]
	switch i.op {
	case "dup":
		return fmt.Sprintf("dup.%s %s, %s", decodeArrangement(i.q, i.size), i.vd, i.gpr)
	case "ins":
		return fmt.Sprintf("mov.%s %s[%d], %s", sz, i.vd, i.idx, i.gpr)
	case "smov":
		return fmt.Sprintf("smov %s, %s.%s[%d]", i.gpr, i.vd, sz, i.idx)
	default:
		// LLVM alias: umov that fills the whole register - mov
		if i.q == 1 && i.size >= 2 {
			return fmt.Sprintf("mov %s, %s.%s[%d]", i.gpr, i.vd, sz, i.idx)
		}

		return fmt.Sprintf("umov %s, %s.%s[%d]", i.gpr, i.vd, sz, i.idx)
	}
}

func (i SimdCopy) Encode(w io.Writer, pc uint64) (int64, error) {
	var opBits uint32
	switch i.op {
	case "dup":
		opBits = 0
	case "ins":
		opBits = 1
	case "smov":
		opBits = 2
	default:
		opBits = 3
	}

	imm5 := uint32(1)<<i.size | i.idx<<(i.size+1)
	// size is inside imm5 (there are no bits 23:22); 0xC00 is the family's
	// fixed bits (b11=1, b10=1), absent from the ctors' enc
	return writeWord(w, i.enc|i.q<<30|0xC00|opBits<<12|imm5<<16|i.gprNum<<5|i.vdNum)
}

func (i SimdCopy) MarshalJSON() ([]byte, error) {
	return i.marshal(i.op, i.ObjDump(disasm.DefaultViewCtx()), "ASIMD",
		map[string]any{"Vd": i.vd, "R": i.gpr, "imm5": uint32(1)<<i.size | i.idx<<(i.size+1)})
}
