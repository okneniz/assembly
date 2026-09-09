package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AsrReg — asr rd, rn, rm.
type AsrReg struct {
	base

	rd, rn, rm string
}

// newAsrReg - the AsrReg constructor: the struct is assembled only
// here (the Builder method and the decoder call it).
func newAsrReg(b base, rd string, rn string, rm string) AsrReg {
	return AsrReg{
		base: b,
		rd:   rd,
		rn:   rn,
		rm:   rm,
	}
}

const AsrRegX uint32 = 0x9A002800

func (i AsrReg) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("asr %s, %s, %s", i.rd, i.rn, i.rm)
}

func (i AsrReg) Encode(w io.Writer) (int64, error) {
	match, err := sfMatch(i.rd, AsrRegX, 0x1A002800)
	if err != nil {
		return 0, fmt.Errorf("asr: %w", err)
	}

	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("asr: %w", err)
	}

	return writeWord(w, match|rd|rn<<5|rm<<16)
}

// AsrReg — asr rd, rn, rm. Register 31 reads as zr (SP/WSP are
// not allowed — use XZR/WZR); the width is shared by all three registers.
func (Builder) AsrReg(rd, rn, rm Reg) (Instr, error) {
	if err := requireClass(rd, "AsrReg", "rd", "register 31 reads as zr — use XZR/WZR",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireClass(rn, "AsrReg", "rn", "register 31 reads as zr — use XZR/WZR",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireClass(rm, "AsrReg", "rm", "register 31 reads as zr — use XZR/WZR",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireWidth("AsrReg", rd, rn, rm); err != nil {
		return nil, err
	}

	return newAsrReg(base{}, rd.name(), rn.name(), rm.name()), nil
}

func decodeAsrReg(w uint32) Instr {
	return newAsrReg(
		newBase(w),
		armRegName(w&0x1f, w>>31&1 == 1),
		armRegName(w>>5&0x1f, w>>31&1 == 1),
		armRegName(w>>16&0x1f, w>>31&1 == 1),
	)
}
