package arm64

import (
	"errors"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// EorShift — eor rd, rn, rm[, shift #imm6].
type EorShift struct {
	base

	rd, rn, rm string
	imm6       uint32
	shift      string
	isf        bool
}

// newEorShift - the EorShift constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newEorShift(b base, rd Reg, rn Reg, rm Reg, imm Imm6, sh Shift) (EorShift, error) {
	err := requireClass(
		rd,
		"EorShift",
		"rd",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return EorShift{}, err
	}

	err = requireClass(
		rn,
		"EorShift",
		"rn",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return EorShift{}, err
	}

	err = requireClass(
		rm,
		"EorShift",
		"rm",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return EorShift{}, err
	}

	err = requireWidth(
		"EorShift",
		rd,
		rn,
		rm,
	)

	if err != nil {
		return EorShift{}, err
	}

	return EorShift{
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
	EorShiftX uint32 = 0xCA000000
	EorShiftW uint32 = 0x4A000000
)

func (i EorShift) ObjDump(_ disasm.ViewCtx) string {
	if i.imm6 != 0 {
		return fmt.Sprintf("eor %s, %s, %s, %s #%d", i.rd, i.rn, i.rm, i.shift, i.imm6)
	}

	return fmt.Sprintf("eor %s, %s, %s", i.rd, i.rn, i.rm)
}

func (i EorShift) Encode(w io.Writer) (int64, error) {
	match := EorShiftX
	if !i.isf {
		match = EorShiftW
	}

	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("eor: %w", err)
	}

	sh, err := shiftNumByName(i.shift)
	if err != nil {
		return 0, fmt.Errorf("eor: %w", err)
	}

	if i.imm6 > 63 {
		return 0, errors.New("eor: imm6 out of range")
	}

	return writeWord(w, match|rd|rn<<5|i.imm6<<10|rm<<16|sh<<22)
}

func (Builder) EorShift(rd, rn, rm Reg, imm Imm6, sh Shift) (Instr, error) {
	return newEorShift(base{}, rd, rn, rm, imm, sh)
}

func decodeEorShift(w uint32) (Instr, error) {
	in, err := newEorShift(
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
