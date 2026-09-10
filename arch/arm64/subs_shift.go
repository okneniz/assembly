package arm64

import (
	"errors"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// SubsShift — subs rd, rn, rm[, shift #imm6]; pseudo: cmp (Rd = zr), negs (Rn = zr).
type SubsShift struct {
	base

	rd, rn, rm string
	imm6       uint32
	shift      string // lsl/lsr/asr/ror
	isf        bool
}

// newSubsShift - the SubsShift constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newSubsShift(b base, rd Reg, rn Reg, rm Reg, imm Imm6, sh Shift) (SubsShift, error) {
	err := requireClass(
		rd,
		"SubsShift",
		"rd",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return SubsShift{}, err
	}

	err = requireClass(
		rn,
		"SubsShift",
		"rn",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return SubsShift{}, err
	}

	err = requireClass(
		rm,
		"SubsShift",
		"rm",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return SubsShift{}, err
	}

	err = requireWidth(
		"SubsShift",
		rd,
		rn,
		rm,
	)

	if err != nil {
		return SubsShift{}, err
	}

	if err = requireShift(rd, "SubsShift", imm, sh); err != nil {
		return SubsShift{}, err
	}

	return SubsShift{
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
	SubsShiftX uint32 = 0xEB000000
	SubsShiftW uint32 = 0x6B000000
)

func (i SubsShift) ObjDump(_ disasm.ViewCtx) string {
	zr := zeroReg(i.rd)
	if i.rd == zr {
		if i.imm6 != 0 {
			return fmt.Sprintf("cmp %s, %s, %s #%d", i.rn, i.rm, i.shift, i.imm6)
		}

		return fmt.Sprintf("cmp %s, %s", i.rn, i.rm)
	}

	if i.rn == zr {
		if i.imm6 != 0 {
			return fmt.Sprintf("negs %s, %s, %s #%d", i.rd, i.rm, i.shift, i.imm6)
		}

		return fmt.Sprintf("negs %s, %s", i.rd, i.rm)
	}

	if i.imm6 != 0 {
		return fmt.Sprintf("subs %s, %s, %s, %s #%d", i.rd, i.rn, i.rm, i.shift, i.imm6)
	}

	return fmt.Sprintf("subs %s, %s, %s", i.rd, i.rn, i.rm)
}

func (i SubsShift) Encode(w io.Writer) (int64, error) {
	match := SubsShiftX
	if !i.isf {
		match = SubsShiftW
	}

	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("subs: %w", err)
	}

	sh, err := shiftNumByName(i.shift)
	if err != nil {
		return 0, fmt.Errorf("subs: %w", err)
	}

	if i.imm6 > 63 {
		return 0, errors.New("subs: imm6 out of range")
	}

	return writeWord(w, match|rd|rn<<5|i.imm6<<10|rm<<16|sh<<22)
}

func (Builder) SubsShift(rd, rn, rm Reg, imm Imm6, sh Shift) (Instr, error) {
	return newSubsShift(base{}, rd, rn, rm, imm, sh)
}

func decodeSubsShift(w uint32) (Instr, error) {
	in, err := newSubsShift(
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
