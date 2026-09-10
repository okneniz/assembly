package arm64

import (
	"errors"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// BicsShift — bics rd, rn, rm[, shift #imm6].
type BicsShift struct {
	base

	rd, rn, rm string
	imm6       uint32
	shift      string
	isf        bool
}

// newBicsShift - the BicsShift constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newBicsShift(b base, rd Reg, rn Reg, rm Reg, imm Imm6, sh Shift) (BicsShift, error) {
	err := requireClass(
		rd,
		"BicsShift",
		"rd",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return BicsShift{}, err
	}

	err = requireClass(
		rn,
		"BicsShift",
		"rn",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return BicsShift{}, err
	}

	err = requireClass(
		rm,
		"BicsShift",
		"rm",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return BicsShift{}, err
	}

	err = requireWidth(
		"BicsShift",
		rd,
		rn,
		rm,
	)

	if err != nil {
		return BicsShift{}, err
	}

	return BicsShift{
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
	BicsShiftX uint32 = 0xEA200000
	BicsShiftW uint32 = 0x6A200000
)

func (i BicsShift) ObjDump(_ disasm.ViewCtx) string {
	if i.imm6 != 0 {
		return fmt.Sprintf("bics %s, %s, %s, %s #%d", i.rd, i.rn, i.rm, i.shift, i.imm6)
	}

	return fmt.Sprintf("bics %s, %s, %s", i.rd, i.rn, i.rm)
}

func (i BicsShift) Encode(w io.Writer) (int64, error) {
	match := BicsShiftX
	if !i.isf {
		match = BicsShiftW
	}

	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("bics: %w", err)
	}

	sh, err := shiftNumByName(i.shift)
	if err != nil {
		return 0, fmt.Errorf("bics: %w", err)
	}

	if i.imm6 > 63 {
		return 0, errors.New("bics: imm6 out of range")
	}

	return writeWord(w, match|rd|rn<<5|i.imm6<<10|rm<<16|sh<<22)
}

func (Builder) BicsShift(rd, rn, rm Reg, imm Imm6, sh Shift) (Instr, error) {
	return newBicsShift(base{}, rd, rn, rm, imm, sh)
}

func decodeBicsShift(w uint32) (Instr, error) {
	in, err := newBicsShift(
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
