package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Fp3 — op rd, rn, rm (fadd/fsub/fmul/fdiv/fmax/fmin; s/d).
type Fp3 struct {
	base

	op         string
	rd, rn, rm string
	enc        uint32
}

func decodeFp3Of(op string, enc uint32, k fpKind) func(uint32, uint64) Instr {
	return func(w uint32, addr uint64) Instr {
		return Fp3{
			base: newBase(addr, w),
			op:   op,
			rd:   fpReg(w&0x1f, k),
			rn:   fpReg(w>>5&0x1f, k),
			rm:   fpReg(w>>16&0x1f, k),
			enc:  enc,
		}
	}
}

func (i Fp3) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("%s %s, %s, %s", i.op, i.rd, i.rn, i.rm)
}

func (i Fp3) Encode(w io.Writer, pc uint64) (int64, error) {
	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", i.op, err)
	}

	return writeWord(w, i.enc|rd|rn<<5|rm<<16)
}

func (i Fp3) MarshalJSON() ([]byte, error) {
	return i.marshal(
		i.op,
		i.ObjDump(disasm.DefaultViewCtx()),
		"FP",
		map[string]any{"Rd": i.rd, "Rn": i.rn, "Rm": i.rm},
	)
}
