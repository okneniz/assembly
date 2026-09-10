package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Fp4 — op rd, rn, rm, ra (fmadd/fnmsub).
type Fp4 struct {
	base

	op             string
	rd, rn, rm, ra string
	enc            uint32
}

func decodeFp4Of(op string, enc uint32) func(uint32) (Instr, error) {
	return func(w uint32) (Instr, error) {
		return Fp4{
			base: newBase(w),
			op:   op,
			rd:   fpReg(w&0x1f, kD),
			rn:   fpReg(w>>5&0x1f, kD),
			rm:   fpReg(w>>16&0x1f, kD),
			ra:   fpReg(w>>10&0x1f, kD),
			enc:  enc,
		}, nil
	}
}

func (i Fp4) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("%s %s, %s, %s, %s", i.op, i.rd, i.rn, i.rm, i.ra)
}

func (i Fp4) Encode(w io.Writer) (int64, error) {
	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", i.op, err)
	}

	ra, err := armRegNum(i.ra)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", i.op, err)
	}

	return writeWord(w, i.enc|rd|rn<<5|ra<<10|rm<<16)
}
