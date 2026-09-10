package arm64

import (
	"errors"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AddsShift — adds rd, rn, rm[, shift #imm6]; pseudo: cmn (Rd = zr).
type AddsShift struct {
	base

	rd, rn, rm string
	imm6       uint32
	shift      string // lsl/lsr/asr/ror
	isf        bool
}

// newAddsShift - the AddsShift constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newAddsShift(b base, rd Reg, rn Reg, rm Reg, imm Imm6, sh Shift) (AddsShift, error) {
	err := requireClass(
		rd,
		"AddsShift",
		"rd",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return AddsShift{}, err
	}

	err = requireClass(
		rn,
		"AddsShift",
		"rn",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return AddsShift{}, err
	}

	err = requireClass(
		rm,
		"AddsShift",
		"rm",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return AddsShift{}, err
	}

	err = requireWidth(
		"AddsShift",
		rd,
		rn,
		rm,
	)

	if err != nil {
		return AddsShift{}, err
	}

	if err = requireShift(rd, "AddsShift", imm, sh); err != nil {
		return AddsShift{}, err
	}

	return AddsShift{
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
	AddsShiftX uint32 = 0xAB000000
	AddsShiftW uint32 = 0x2B000000
)

func (i AddsShift) ObjDump(_ disasm.ViewCtx) string {
	zr := zeroReg(i.rd)
	if i.rd == zr {
		if i.imm6 != 0 {
			return fmt.Sprintf("cmn %s, %s, %s #%d", i.rn, i.rm, i.shift, i.imm6)
		}

		return fmt.Sprintf("cmn %s, %s", i.rn, i.rm)
	}

	if i.imm6 != 0 {
		return fmt.Sprintf("adds %s, %s, %s, %s #%d", i.rd, i.rn, i.rm, i.shift, i.imm6)
	}

	return fmt.Sprintf("adds %s, %s, %s", i.rd, i.rn, i.rm)
}

func (i AddsShift) Encode(w io.Writer) (int64, error) {
	match := AddsShiftX
	if !i.isf {
		match = AddsShiftW
	}

	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("adds: %w", err)
	}

	sh, err := shiftNumByName(i.shift)
	if err != nil {
		return 0, fmt.Errorf("adds: %w", err)
	}

	if i.imm6 > 63 {
		return 0, errors.New("adds: imm6 out of range")
	}

	return writeWord(w, match|rd|rn<<5|i.imm6<<10|rm<<16|sh<<22)
}

func (Builder) AddsShift(rd, rn, rm Reg, imm Imm6, sh Shift) (Instr, error) {
	return newAddsShift(base{}, rd, rn, rm, imm, sh)
}

func decodeAddsShift(w uint32) (Instr, error) {
	in, err := newAddsShift(
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
