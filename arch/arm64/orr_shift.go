package arm64

import (
	"errors"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// OrrShift — orr rd, rn, rm[, shift #imm6]; pseudo: mov (Rn = zr).
type OrrShift struct {
	base

	rd, rn, rm string
	imm6       uint32
	shift      string
	isf        bool
}

// newOrrShift - the OrrShift constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newOrrShift(b base, rd Reg, rn Reg, rm Reg, imm Imm6, sh Shift) (OrrShift, error) {
	err := requireClass(
		rd,
		"OrrShift",
		"rd",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return OrrShift{}, err
	}

	err = requireClass(
		rn,
		"OrrShift",
		"rn",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return OrrShift{}, err
	}

	err = requireClass(
		rm,
		"OrrShift",
		"rm",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return OrrShift{}, err
	}

	err = requireWidth(
		"OrrShift",
		rd,
		rn,
		rm,
	)

	if err != nil {
		return OrrShift{}, err
	}

	return OrrShift{
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
	OrrShiftX uint32 = 0xAA000000
	OrrShiftW uint32 = 0x2A000000
)

func (i OrrShift) ObjDump(_ disasm.ViewCtx) string {
	zr := zeroReg(i.rd)
	if i.rn == zr {
		if i.imm6 != 0 {
			return fmt.Sprintf("mov %s, %s, %s #%d", i.rd, i.rm, i.shift, i.imm6)
		}

		return fmt.Sprintf("mov %s, %s", i.rd, i.rm)
	}

	if i.imm6 != 0 {
		return fmt.Sprintf("orr %s, %s, %s, %s #%d", i.rd, i.rn, i.rm, i.shift, i.imm6)
	}

	return fmt.Sprintf("orr %s, %s, %s", i.rd, i.rn, i.rm)
}

func (i OrrShift) Encode(w io.Writer) (int64, error) {
	match := OrrShiftX
	if !i.isf {
		match = OrrShiftW
	}

	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("orr: %w", err)
	}

	sh, err := shiftNumByName(i.shift)
	if err != nil {
		return 0, fmt.Errorf("orr: %w", err)
	}

	if i.imm6 > 63 {
		return 0, errors.New("orr: imm6 out of range")
	}

	return writeWord(w, match|rd|rn<<5|i.imm6<<10|rm<<16|sh<<22)
}

func (Builder) OrrShift(rd, rn, rm Reg, imm Imm6, sh Shift) (Instr, error) {
	return newOrrShift(base{}, rd, rn, rm, imm, sh)
}

func decodeOrrShift(w uint32) (Instr, error) {
	in, err := newOrrShift(
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
