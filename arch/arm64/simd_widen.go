package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// SimdWiden — the widening three-same family: saddw/ssubw/uaddw/usubw.
// U (bit 29) selects s/u, bit 12 - add/sub, Q - the "2" form. The result
// and Rn have the arrangement (size+1), the Rm source - size (narrow):
// saddw=0x0E201000, ssubw=0x0E203000, uaddw=0x2E201000, usubw=0x2E203000,
// mask 0xBFA0FC00.
type SimdWiden struct {
	base

	op       string
	q, size  uint32
	rd, rn   string
	rm       string
	enc      uint32
	rdN, rnN uint32
	rmN      uint32
}

func decodeSimdWidenOf(op string, enc uint32) func(w uint32, addr uint64) Instr {
	return func(w uint32, addr uint64) Instr {
		name := op
		if w>>30&1 == 1 {
			name += "2"
		}

		return SimdWiden{
			base: newBase(addr, w),
			op:   name,
			q:    w >> 30 & 1,
			size: w >> 22 & 3,
			rd:   vReg(w & 0x1f),
			rn:   vReg(w >> 5 & 0x1f),
			rm:   vReg(w >> 16 & 0x1f),
			enc:  enc,
			rdN:  w & 0x1f,
			rnN:  w >> 5 & 0x1f,
			rmN:  w >> 16 & 0x1f,
		}
	}
}

func (i SimdWiden) ObjDump(_ disasm.ViewCtx) string {
	arr := decodeArrangement(i.q, i.size+1)
	return fmt.Sprintf("%s.%s %s, %s, %s", i.op, arr, i.rd, i.rn, i.rm)
}

func (i SimdWiden) Encode(w io.Writer, pc uint64) (int64, error) {
	return writeWord(w, i.enc|i.q<<30|i.size<<22|i.rmN<<16|i.rnN<<5|i.rdN)
}

func (i SimdWiden) MarshalJSON() ([]byte, error) {
	return i.marshal(i.op, i.ObjDump(disasm.DefaultViewCtx()), "ASIMD",
		map[string]any{"Rd": i.rd, "Rn": i.rn, "Rm": i.rm})
}
