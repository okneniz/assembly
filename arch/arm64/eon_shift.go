package arm64

import (
	"errors"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// EonShift — eon rd, rn, rm[, shift #imm6].
type EonShift struct {
	base

	rd, rn, rm string
	imm6       uint32
	shift      string
	isf        bool
}

// newEonShift - the EonShift constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newEonShift(b base, rd Reg, rn Reg, rm Reg, imm Imm6, sh Shift) (EonShift, error) {
	err := requireClass(
		rd,
		"EonShift",
		"rd",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return EonShift{}, err
	}

	err = requireClass(
		rn,
		"EonShift",
		"rn",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return EonShift{}, err
	}

	err = requireClass(
		rm,
		"EonShift",
		"rm",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return EonShift{}, err
	}

	err = requireWidth(
		"EonShift",
		rd,
		rn,
		rm,
	)

	if err != nil {
		return EonShift{}, err
	}

	return EonShift{
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
	EonShiftX uint32 = 0xCA200000
	EonShiftW uint32 = 0x4A200000
)

func (i EonShift) ObjDump(_ disasm.ViewCtx) string {
	if i.imm6 != 0 {
		return fmt.Sprintf("eon %s, %s, %s, %s #%d", i.rd, i.rn, i.rm, i.shift, i.imm6)
	}

	return fmt.Sprintf("eon %s, %s, %s", i.rd, i.rn, i.rm)
}

func (i EonShift) Encode(w io.Writer) (int64, error) {
	match := EonShiftX
	if !i.isf {
		match = EonShiftW
	}

	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("eon: %w", err)
	}

	sh, err := shiftNumByName(i.shift)
	if err != nil {
		return 0, fmt.Errorf("eon: %w", err)
	}

	if i.imm6 > 63 {
		return 0, errors.New("eon: imm6 out of range")
	}

	return writeWord(w, match|rd|rn<<5|i.imm6<<10|rm<<16|sh<<22)
}

func (Builder) EonShift(rd, rn, rm Reg, imm Imm6, sh Shift) (Instr, error) {
	return newEonShift(base{}, rd, rn, rm, imm, sh)
}

func decodeEonShift(w uint32) (Instr, error) {
	in, err := newEonShift(
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
