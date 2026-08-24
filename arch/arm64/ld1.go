package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Ld1 — ld1.Arr { vN, ... }{[idx]}, [rn]{, #imm | rm} (structural loads).
type Ld1 struct {
	base

	regList string // "{ v0, v1 }" | "{ v0 }[2]"
	rn      string
	name    string
	arr     string
	postReg string
	postImm uint32
	hasPost bool
	enc     uint32
	rtNum   uint32
	count   int
	opcode  uint32
	size    uint32
	q       uint32
	isElem  bool
}

func decodeLd1Of(enc uint32) func(w uint32, addr uint64) Instr {
	return func(w uint32, addr uint64) Instr {
		opcode, size, q := w>>12&0xf, w>>10&3, w>>30&1
		name, arr, count, isElem := ldStructDecode(opcode, size, q, w>>22&1)
		rt := w & 0x1f
		list := "{ " + regListStr(rt, count) + " }"
		if w>>24&1 == 1 {
			// single-structure (bits[29:24]=001101): the layout is DIFFERENT,
			// not like multi — size lives in bits[15:14], the element index =
			// Q with the truncated bits [13:10], bits[15:14]=11 → ld1r.
			name, arr, list, size = ldElemDecode(w, q, rt)
		}

		i := Ld1{
			base:    newBase(addr, w),
			regList: list,
			rn:      regNameXSP(w >> 5 & 0x1f),
			name:    name,
			arr:     arr,
			enc:     enc,
			rtNum:   rt,
			count:   count,
			opcode:  opcode,
			size:    size,
			q:       q,
			isElem:  isElem,
		}
		// post-index: bit23 (Rm=11111 → immediate, otherwise post-index by
		// register — objdump prints the register name); bit22 = L.
		if w>>23&1 == 1 {
			if rm := w >> 16 & 0x1f; rm != 0x1f {
				i.hasPost = true
				i.postReg = regNameX(rm)
				return i
			}

			elemBytes := uint32(1) << size
			regBytes := uint32(8)
			if q == 1 {
				regBytes = 16
			}

			postImm := regBytes * uint32(count)
			if opcode == 0xe {
				postImm = elemBytes * 4
			}

			if opcode == 0xc {
				postImm = elemBytes
			}

			if w>>24&1 == 1 {
				postImm = elemBytes
				if name == "ld1r" {
					postImm = regBytes // ld1r writes the whole register
				}

				if name == "ld4r" {
					postImm = regBytes // ld4r: post = one register (objdump: #16 for .4s)
				}
			}

			i.hasPost, i.postImm = true, postImm
		}

		return i
	}
}

// ldElemDecode — single-structure forms (element and replicate). The
// ld1r/ld4r split — by bit 21.
func ldElemDecode(w, q, rt uint32) (name, arr, list string, size uint32) {
	if w>>14&3 == 3 {
		sz := w >> 10 & 3
		if w>>21&1 == 1 {
			return "ld4r", decodeArrangement(q, sz), "{ " + regListStr(rt, 4) + " }", sz
		}

		return "ld1r", decodeArrangement(q, sz), "{ " + regListStr(rt, 1) + " }", sz
	}

	switch w >> 14 & 3 {
	case 0: // .b: index = Q:bits[12:10]
		idx := q<<3 | w>>10&7
		return "ld1", "b", fmt.Sprintf("{ %s }[%d]", vReg(rt), idx), 0
	case 1: // .h: index = Q:bits[12:11]
		idx := q<<2 | w>>11&3
		return "ld1", "h", fmt.Sprintf("{ %s }[%d]", vReg(rt), idx), 1
	case 2:
		if w>>10&1 == 0 { // .s: index = Q:bit12
			idx := q<<1 | w>>12&1
			return "ld1", "s", fmt.Sprintf("{ %s }[%d]", vReg(rt), idx), 2
		} // .d: index = Q

		return "ld1", "d", fmt.Sprintf("{ %s }[%d]", vReg(rt), q), 3
	}

	return "ld1", "b", "", 0 // unreachable
}

func (i Ld1) ObjDump(_ disasm.ViewCtx) string {
	s := fmt.Sprintf("%s.%s %s, [%s]", i.name, i.arr, i.regList, i.rn)
	if i.hasPost {
		if i.postReg != "" {
			s += ", " + i.postReg
		} else {
			s += fmt.Sprintf(", #%d", i.postImm)
		}
	}

	return s
}

func (i Ld1) Encode(w io.Writer, pc uint64) (int64, error) {
	rn, err := armRegNum(i.rn)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", i.name, err)
	}

	rm := uint32(0x1f)
	if i.postReg != "" {
		rm, err = armRegNum(i.postReg)
		if err != nil {
			return 0, fmt.Errorf("%s: %w", i.name, err)
		}
	}

	return writeWord(w, i.enc|i.rtNum|rn<<5|rm<<16)
}

func (i Ld1) MarshalJSON() ([]byte, error) {
	return i.marshal(
		i.name,
		i.ObjDump(disasm.DefaultViewCtx()),
		"ASIMD",
		map[string]any{"Rt": i.rtNum, "Rn": i.rn},
	)
}
