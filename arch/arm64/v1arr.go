package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// V1arr — op.16b vd, vn (aese/aesmc).
type V1arr struct {
	base

	op     string
	rd, rn string
	enc    uint32
}

func decodeV1arrOf(op string, enc uint32) func(uint32) (Instr, error) {
	return func(w uint32) (Instr, error) {
		return V1arr{
			base: newBase(w),
			op:   op,
			rd:   vReg(w & 0x1f),
			rn:   vReg(w >> 5 & 0x1f),
			enc:  enc,
		}, nil
	}
}

func (i V1arr) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("%s.16b %s, %s", i.op, i.rd, i.rn)
}

func (i V1arr) Encode(w io.Writer) (int64, error) {
	rd, rn, err := regNums2(i.rd, i.rn)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", i.op, err)
	}

	return writeWord(w, i.enc|rd|rn<<5)
}
