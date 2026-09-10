package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Sdiv — sdiv rd, rn, rm.
type Sdiv struct {
	base

	rd, rn, rm string
}

// newSdiv - the Sdiv constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newSdiv(b base, rd Reg, rn Reg, rm Reg) (Sdiv, error) {
	err := requireClass(
		rd,
		"Sdiv",
		"rd",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return Sdiv{}, err
	}

	err = requireClass(
		rn,
		"Sdiv",
		"rn",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return Sdiv{}, err
	}

	err = requireClass(
		rm,
		"Sdiv",
		"rm",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return Sdiv{}, err
	}

	err = requireWidth(
		"Sdiv",
		rd,
		rn,
		rm,
	)

	if err != nil {
		return Sdiv{}, err
	}

	return Sdiv{
		base: b,
		rd:   rd.name(),
		rn:   rn.name(),
		rm:   rm.name(),
	}, nil
}

const SdivX uint32 = 0x9AC00C00

func (i Sdiv) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("sdiv %s, %s, %s", i.rd, i.rn, i.rm)
}

func (i Sdiv) Encode(w io.Writer) (int64, error) {
	match, err := sfMatch(i.rd, SdivX, 0x1AC00C00)
	if err != nil {
		return 0, fmt.Errorf("sdiv: %w", err)
	}

	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("sdiv: %w", err)
	}

	return writeWord(w, match|rd|rn<<5|rm<<16)
}

func (Builder) Sdiv(rd, rn, rm Reg) (Instr, error) {
	return newSdiv(base{}, rd, rn, rm)
}

func decodeSdiv(w uint32) (Instr, error) {
	in, err := newSdiv(
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
