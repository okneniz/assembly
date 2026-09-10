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

// newLslReg - the LslReg constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newLslReg(b base, rd Reg, rn Reg, rm Reg) (LslReg, error) {
	err := requireClass(
		rd,
		"LslReg",
		"rd",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return LslReg{}, err
	}

	err = requireClass(
		rn,
		"LslReg",
		"rn",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return LslReg{}, err
	}

	err = requireClass(
		rm,
		"LslReg",
		"rm",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return LslReg{}, err
	}

	err = requireWidth(
		"LslReg",
		rd,
		rn,
		rm,
	)

	if err != nil {
		return LslReg{}, err
	}

	return LslReg{
		base: b,
		rd:   rd.name(),
		rn:   rn.name(),
		rm:   rm.name(),
	}, nil
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

func (Builder) LslReg(rd, rn, rm Reg) (Instr, error) {
	return newLslReg(base{}, rd, rn, rm)
}

func decodeLslReg(w uint32) (Instr, error) {
	in, err := newLslReg(
		newBase(w),
		gprOf(w&0x1f, w>>31&1 == 1),
		gprOf(w>>5&0x1f, w>>31&1 == 1),
		gprOf(w>>16&0x1f, w>>31&1 == 1),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}
