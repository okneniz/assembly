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

// newAsrReg - the AsrReg constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newAsrReg(b base, rd Reg, rn Reg, rm Reg) (AsrReg, error) {
	err := requireClass(
		rd,
		"AsrReg",
		"rd",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return AsrReg{}, err
	}

	err = requireClass(
		rn,
		"AsrReg",
		"rn",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return AsrReg{}, err
	}

	err = requireClass(
		rm,
		"AsrReg",
		"rm",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return AsrReg{}, err
	}

	err = requireWidth(
		"AsrReg",
		rd,
		rn,
		rm,
	)

	if err != nil {
		return AsrReg{}, err
	}

	return AsrReg{
		base: b,
		rd:   rd.name(),
		rn:   rn.name(),
		rm:   rm.name(),
	}, nil
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

func (Builder) AsrReg(rd, rn, rm Reg) (Instr, error) {
	return newAsrReg(base{}, rd, rn, rm)
}

func decodeAsrReg(w uint32) (Instr, error) {
	in, err := newAsrReg(
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
