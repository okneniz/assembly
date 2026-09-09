package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// LslReg — lsl rd, rn, rm.
type LslReg struct {
	base

	rd, rn, rm string
}

// newLslReg - the LslReg constructor: the struct is assembled only
// here (the Builder method and the decoder call it).
func newLslReg(b base, rd string, rn string, rm string) LslReg {
	return LslReg{
		base: b,
		rd:   rd,
		rn:   rn,
		rm:   rm,
	}
}

const LslRegX uint32 = 0x9A002000

func (i LslReg) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("lsl %s, %s, %s", i.rd, i.rn, i.rm)
}

func (i LslReg) Encode(w io.Writer) (int64, error) {
	match, err := sfMatch(i.rd, LslRegX, 0x1A002000)
	if err != nil {
		return 0, fmt.Errorf("lsl: %w", err)
	}

	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("lsl: %w", err)
	}

	return writeWord(w, match|rd|rn<<5|rm<<16)
}

// LslReg — lsl rd, rn, rm. Register 31 reads as zr (SP/WSP are
// not allowed — use XZR/WZR); the width is shared by all three registers.
func (Builder) LslReg(rd, rn, rm Reg) (Instr, error) {
	if err := requireClass(rd, "LslReg", "rd", "register 31 reads as zr — use XZR/WZR",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireClass(rn, "LslReg", "rn", "register 31 reads as zr — use XZR/WZR",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireClass(rm, "LslReg", "rm", "register 31 reads as zr — use XZR/WZR",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireWidth("LslReg", rd, rn, rm); err != nil {
		return nil, err
	}

	return newLslReg(base{}, rd.name(), rn.name(), rm.name()), nil
}

func decodeLslReg(w uint32) Instr {
	return newLslReg(
		newBase(w),
		armRegName(w&0x1f, w>>31&1 == 1),
		armRegName(w>>5&0x1f, w>>31&1 == 1),
		armRegName(w>>16&0x1f, w>>31&1 == 1),
	)
}
