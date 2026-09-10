package arm64

import (
	"errors"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// SubShift — sub rd, rn, rm[, shift #imm6]; pseudo: neg (Rn = zr).
type SubShift struct {
	base

	rd, rn, rm string
	imm6       uint32
	shift      string // lsl/lsr/asr/ror
	isf        bool
}

// newSubShift - the SubShift constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newSubShift(b base, rd Reg, rn Reg, rm Reg, imm Imm6, sh Shift) (SubShift, error) {
	for _, r := range []struct {
		reg Reg
		op  string
	}{{
		rd,
		"rd",
	}, {
		rn,
		"rn",
	}, {
		rm,
		"rm",
	}} {
		err := requireClass(
			r.reg,
			"SubShift",
			r.op,
			"register 31 reads as zr - use XZR/WZR",
			classX,
			classW,
			classXZR,
			classWZR,
		)

		if err != nil {
			return SubShift{}, err
		}
	}

	err := requireWidth(
		"SubShift",
		rd,
		rn,
		rm,
	)

	if err != nil {
		return SubShift{}, err
	}

	if err = requireShift(rd, "SubShift", imm, sh); err != nil {
		return SubShift{}, err
	}

	return SubShift{
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
	SubShiftX uint32 = 0xCB000000
	SubShiftW uint32 = 0x4B000000
)

func (i SubShift) ObjDump(_ disasm.ViewCtx) string {
	zr := zeroReg(i.rd)
	if i.rn == zr {
		if i.imm6 != 0 {
			return fmt.Sprintf("neg %s, %s, %s #%d", i.rd, i.rm, i.shift, i.imm6)
		}

		return fmt.Sprintf("neg %s, %s", i.rd, i.rm)
	}

	if i.imm6 != 0 {
		return fmt.Sprintf("sub %s, %s, %s, %s #%d", i.rd, i.rn, i.rm, i.shift, i.imm6)
	}

	return fmt.Sprintf("sub %s, %s, %s", i.rd, i.rn, i.rm)
}

func (i SubShift) Encode(w io.Writer) (int64, error) {
	match := SubShiftX
	if !i.isf {
		match = SubShiftW
	}

	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("sub: %w", err)
	}

	sh, err := shiftNumByName(i.shift)
	if err != nil {
		return 0, fmt.Errorf("sub: %w", err)
	}

	if i.imm6 > 63 {
		return 0, errors.New("sub: imm6 out of range")
	}

	return writeWord(w, match|rd|rn<<5|i.imm6<<10|rm<<16|sh<<22)
}

func (Builder) SubShift(rd, rn, rm Reg, imm Imm6, sh Shift) (Instr, error) {
	return newSubShift(base{}, rd, rn, rm, imm, sh)
}

func decodeSubShift(w uint32) (Instr, error) {
	in, err := newSubShift(
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
