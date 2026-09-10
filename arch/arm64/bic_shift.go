package arm64

import (
	"errors"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// BicShift — bic rd, rn, rm[, shift #imm6].
type BicShift struct {
	base

	rd, rn, rm string
	imm6       uint32
	shift      string
	isf        bool
}

// newBicShift - the BicShift constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newBicShift(b base, rd Reg, rn Reg, rm Reg, imm Imm6, sh Shift) (BicShift, error) {
	err := requireClass(
		rd,
		"BicShift",
		"rd",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return BicShift{}, err
	}

	err = requireClass(
		rn,
		"BicShift",
		"rn",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return BicShift{}, err
	}

	err = requireClass(
		rm,
		"BicShift",
		"rm",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return BicShift{}, err
	}

	err = requireWidth(
		"BicShift",
		rd,
		rn,
		rm,
	)

	if err != nil {
		return BicShift{}, err
	}

	return BicShift{
		base:  b,
		rd:    rd.name(),
		rn:    rn.name(),
		rm:    rm.name(),
		imm6:  imm.v,
		shift: sh.String(),
		isf:   rd.Is64(),
	}, nil
}

const (
	BicShiftX uint32 = 0x8A200000
	BicShiftW uint32 = 0x0A200000
)

func (i BicShift) ObjDump(_ disasm.ViewCtx) string {
	if i.imm6 != 0 {
		return fmt.Sprintf("bic %s, %s, %s, %s #%d", i.rd, i.rn, i.rm, i.shift, i.imm6)
	}

	return fmt.Sprintf("bic %s, %s, %s", i.rd, i.rn, i.rm)
}

func (i BicShift) Encode(w io.Writer) (int64, error) {
	match := BicShiftX
	if !i.isf {
		match = BicShiftW
	}

	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("bic: %w", err)
	}

	sh, err := shiftNumByName(i.shift)
	if err != nil {
		return 0, fmt.Errorf("bic: %w", err)
	}

	if i.imm6 > 63 {
		return 0, errors.New("bic: imm6 out of range")
	}

	return writeWord(w, match|rd|rn<<5|i.imm6<<10|rm<<16|sh<<22)
}

func (Builder) BicShift(rd, rn, rm Reg, imm Imm6, sh Shift) (Instr, error) {
	return newBicShift(base{}, rd, rn, rm, imm, sh)
}

func decodeBicShift(w uint32) (Instr, error) {
	in, err := newBicShift(
		newBase(w),
		gprOf(w&0x1f, w>>31&1 == 1),
		gprOf(w>>5&0x1f, w>>31&1 == 1),
		gprOf(w>>16&0x1f, w>>31&1 == 1),
		imm6Of(w>>10&0x3f),
		shiftOf(w>>22&3),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}
