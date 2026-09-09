package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// LsrReg — lsr rd, rn, rm.
type LsrReg struct {
	base

	rd, rn, rm string
}

// newLsrReg - the LsrReg constructor: the struct is assembled only
// here (the Builder method and the decoder call it).
func newLsrReg(b base, rd string, rn string, rm string) LsrReg {
	return LsrReg{
		base: b,
		rd:   rd,
		rn:   rn,
		rm:   rm,
	}
}

const LsrRegX uint32 = 0x9A002400

func (i LsrReg) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("lsr %s, %s, %s", i.rd, i.rn, i.rm)
}

func (i LsrReg) Encode(w io.Writer) (int64, error) {
	match, err := sfMatch(i.rd, LsrRegX, 0x1A002400)
	if err != nil {
		return 0, fmt.Errorf("lsr: %w", err)
	}

	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("lsr: %w", err)
	}

	return writeWord(w, match|rd|rn<<5|rm<<16)
}

// LsrReg — lsr rd, rn, rm. Register 31 reads as zr (SP/WSP are
// not allowed — use XZR/WZR); the width is shared by all three registers.
func (Builder) LsrReg(rd, rn, rm Reg) (Instr, error) {
	if err := requireClass(rd, "LsrReg", "rd", "register 31 reads as zr — use XZR/WZR",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireClass(rn, "LsrReg", "rn", "register 31 reads as zr — use XZR/WZR",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireClass(rm, "LsrReg", "rm", "register 31 reads as zr — use XZR/WZR",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireWidth("LsrReg", rd, rn, rm); err != nil {
		return nil, err
	}

	return newLsrReg(base{}, rd.name(), rn.name(), rm.name()), nil
}

func decodeLsrReg(w uint32) Instr {
	return newLsrReg(
		newBase(w),
		armRegName(w&0x1f, w>>31&1 == 1),
		armRegName(w>>5&0x1f, w>>31&1 == 1),
		armRegName(w>>16&0x1f, w>>31&1 == 1),
	)
}
