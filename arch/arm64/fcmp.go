package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Fcmp — fcmp rn, rm | fcmp rn, #0.0 (Rm=0 encodes #0.0).
type Fcmp struct {
	base

	rn, rm string
	withRM bool
	enc0   uint32 // base without the Rm bit (#0.0: opcode2=00)
	encR   uint32 // base of the register form
	k      fpKind
}

func decodeFcmpOf(withRM bool, enc uint32, k fpKind) func(uint32, uint64) Instr {
	return func(w uint32, addr uint64) Instr {
		return Fcmp{
			base:   newBase(addr, w),
			rn:     fpReg(w>>5&0x1f, k),
			rm:     fpReg(w>>16&0x1f, k),
			withRM: withRM,
			enc0:   enc,
			encR:   enc,
			k:      k,
		}
	}
}

func (i Fcmp) ObjDump(_ disasm.ViewCtx) string {
	if !i.withRM {
		return fmt.Sprintf("fcmp %s, #0.0", i.rn)
	}

	return fmt.Sprintf("fcmp %s, %s", i.rn, i.rm)
}

func (i Fcmp) Encode(w io.Writer, pc uint64) (int64, error) {
	rn, err := armRegNum(i.rn)
	if err != nil {
		return 0, fmt.Errorf("fcmp: %w", err)
	}

	word := i.enc0 | rn<<5
	if i.withRM {
		rm, err := armRegNum(i.rm)
		if err != nil {
			return 0, fmt.Errorf("fcmp: %w", err)
		}

		word = i.encR | rn<<5 | rm<<16
	}

	return writeWord(w, word)
}

func (i Fcmp) MarshalJSON() ([]byte, error) {
	return i.marshal("fcmp", i.ObjDump(disasm.DefaultViewCtx()), "FP", map[string]any{"Rn": i.rn})
}
