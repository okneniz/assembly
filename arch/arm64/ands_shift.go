package arm64

import (
	"errors"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AndsShift — ands rd, rn, rm[, shift #imm6]; pseudo: tst (Rd = zr, imm6 = 0).
type AndsShift struct {
	base

	rd, rn, rm string
	imm6       uint32
	shift      string
	isf        bool
}

// newAndsShift - the AndsShift constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newAndsShift(b base, rd Reg, rn Reg, rm Reg, imm Imm6, sh Shift) (AndsShift, error) {
	err := requireClass(
		rd,
		"AndsShift",
		"rd",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return AndsShift{}, err
	}

	err = requireClass(
		rn,
		"AndsShift",
		"rn",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return AndsShift{}, err
	}

	err = requireClass(
		rm,
		"AndsShift",
		"rm",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return AndsShift{}, err
	}

	err = requireWidth(
		"AndsShift",
		rd,
		rn,
		rm,
	)

	if err != nil {
		return AndsShift{}, err
	}

	return AndsShift{
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
	AndsShiftX uint32 = 0xEA000000
	AndsShiftW uint32 = 0x6A000000
)

func (i AndsShift) ObjDump(_ disasm.ViewCtx) string {
	zr := zeroReg(i.rd)
	if i.rd == zr && i.imm6 == 0 {
		return fmt.Sprintf("tst %s, %s", i.rn, i.rm)
	}

	if i.imm6 != 0 {
		return fmt.Sprintf("ands %s, %s, %s, %s #%d", i.rd, i.rn, i.rm, i.shift, i.imm6)
	}

	return fmt.Sprintf("ands %s, %s, %s", i.rd, i.rn, i.rm)
}

func (i AndsShift) Encode(w io.Writer) (int64, error) {
	match := AndsShiftX
	if !i.isf {
		match = AndsShiftW
	}

	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("ands: %w", err)
	}

	sh, err := shiftNumByName(i.shift)
	if err != nil {
		return 0, fmt.Errorf("ands: %w", err)
	}

	if i.imm6 > 63 {
		return 0, errors.New("ands: imm6 out of range")
	}

	return writeWord(w, match|rd|rn<<5|i.imm6<<10|rm<<16|sh<<22)
}

func (Builder) AndsShift(rd, rn, rm Reg, imm Imm6, sh Shift) (Instr, error) {
	return newAndsShift(base{}, rd, rn, rm, imm, sh)
}

func decodeAndsShift(w uint32) (Instr, error) {
	in, err := newAndsShift(
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
