package arm64

import (
	"errors"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AddShift — add rd, rn, rm[, shift #imm6].
type AddShift struct {
	base

	rd, rn, rm string
	imm6       uint32
	shift      string // lsl/lsr/asr/ror
	isf        bool
}

// newAddShift - the AddShift constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newAddShift(b base, rd Reg, rn Reg, rm Reg, imm Imm6, sh Shift) (AddShift, error) {
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
			"AddShift",
			r.op,
			"register 31 reads as zr — use XZR/WZR",
			classX,
			classW,
			classXZR,
			classWZR,
		)

		if err != nil {
			return AddShift{}, err
		}
	}

	err := requireWidth(
		"AddShift",
		rd,
		rn,
		rm,
	)

	if err != nil {
		return AddShift{}, err
	}

	if err = requireShift(rd, "AddShift", imm, sh); err != nil {
		return AddShift{}, err
	}

	return AddShift{
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
	AddShiftX uint32 = 0x8B000000
	AddShiftW uint32 = 0x0B000000
)

func (i AddShift) ObjDump(_ disasm.ViewCtx) string {
	if i.imm6 != 0 {
		return fmt.Sprintf("add %s, %s, %s, %s #%d", i.rd, i.rn, i.rm, i.shift, i.imm6)
	}

	return fmt.Sprintf("add %s, %s, %s", i.rd, i.rn, i.rm)
}

func (i AddShift) Encode(w io.Writer) (int64, error) {
	match := AddShiftX
	if !i.isf {
		match = AddShiftW
	}

	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("add: %w", err)
	}

	sh, err := shiftNumByName(i.shift)
	if err != nil {
		return 0, fmt.Errorf("add: %w", err)
	}

	if i.imm6 > 63 {
		return 0, errors.New("add: imm6 out of range")
	}

	return writeWord(w, match|rd|rn<<5|i.imm6<<10|rm<<16|sh<<22)
}

func (Builder) AddShift(rd, rn, rm Reg, imm Imm6, sh Shift) (Instr, error) {
	return newAddShift(base{}, rd, rn, rm, imm, sh)
}

func decodeAddShift(w uint32) (Instr, error) {
	in, err := newAddShift(
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
