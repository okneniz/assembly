package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Simd2 - op.Arr vd, vn (a single vector source operand).
type Simd2 struct {
	base

	op      string
	rd, rn  string
	arr     string
	enc     uint32
	q, size uint32
}

func decodeSimd2Of(op string, enc uint32) func(uint32, uint64) Instr {
	return func(w uint32, addr uint64) Instr {
		q, size := w>>30&1, w>>22&3
		return Simd2{
			base: newBase(addr, w),
			op:   op,
			rd:   vReg(w & 0x1f),
			rn:   vReg(w >> 5 & 0x1f),
			arr:  decodeArrangement(q, size),
			enc:  enc,
			q:    q,
			size: size,
		}
	}
}

func (i Simd2) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("%s.%s %s, %s", i.op, i.arr, i.rd, i.rn)
}

func (i Simd2) Encode(w io.Writer, pc uint64) (int64, error) {
	rd, rn, err := regNums2(i.rd, i.rn)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", i.op, err)
	}

	// the enc constants are in the .8b form (Q=0, size=0); q/size come from the arrangement
	return writeWord(w, i.enc|i.q<<30|i.size<<22|rd|rn<<5)
}

func (i Simd2) MarshalJSON() ([]byte, error) {
	return i.marshal(
		i.op,
		i.ObjDump(disasm.DefaultViewCtx()),
		"ASIMD",
		map[string]any{"Rd": i.rd, "Rn": i.rn},
	)
}
