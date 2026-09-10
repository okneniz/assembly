package arm64

import (
	"errors"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AndShift — and rd, rn, rm[, shift #imm6].
type AndShift struct {
	base

	rd, rn, rm string
	imm6       uint32
	shift      string
	isf        bool
}

// newAndShift - the AndShift constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newAndShift(b base, rd Reg, rn Reg, rm Reg, imm Imm6, sh Shift) (AndShift, error) {
	err := requireClass(
		rd,
		"AndShift",
		"rd",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return AndShift{}, err
	}

	err = requireClass(
		rn,
		"AndShift",
		"rn",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return AndShift{}, err
	}

	err = requireClass(
		rm,
		"AndShift",
		"rm",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return AndShift{}, err
	}

	err = requireWidth(
		"AndShift",
		rd,
		rn,
		rm,
	)

	if err != nil {
		return AndShift{}, err
	}

	return AndShift{
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
	AndShiftX uint32 = 0x8A000000
	AndShiftW uint32 = 0x0A000000
)

func (i AndShift) ObjDump(_ disasm.ViewCtx) string {
	if i.imm6 != 0 {
		return fmt.Sprintf("and %s, %s, %s, %s #%d", i.rd, i.rn, i.rm, i.shift, i.imm6)
	}

	return fmt.Sprintf("and %s, %s, %s", i.rd, i.rn, i.rm)
}

func (i AndShift) Encode(w io.Writer) (int64, error) {
	match := AndShiftX
	if !i.isf {
		match = AndShiftW
	}

	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("and: %w", err)
	}

	sh, err := shiftNumByName(i.shift)
	if err != nil {
		return 0, fmt.Errorf("and: %w", err)
	}

	if i.imm6 > 63 {
		return 0, errors.New("and: imm6 out of range")
	}

	return writeWord(w, match|rd|rn<<5|i.imm6<<10|rm<<16|sh<<22)
}

func (Builder) AndShift(rd, rn, rm Reg, imm Imm6, sh Shift) (Instr, error) {
	return newAndShift(base{}, rd, rn, rm, imm, sh)
}

func decodeAndShift(w uint32) (Instr, error) {
	in, err := newAndShift(
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
