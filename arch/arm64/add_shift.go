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

// newAddShift - the AddShift constructor: the struct is assembled only
// here (the Builder method and the decoder call it).
func newAddShift(
	b base,
	rd string,
	rn string,
	rm string,
	imm6 uint32,
	shift string,
	isf bool,
) AddShift {
	return AddShift{
		base:  b,
		rd:    rd,
		rn:    rn,
		rm:    rm,
		imm6:  imm6,
		shift: shift,
		isf:   isf,
	}
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

// AddShift — add rd, rn, rm[, shift #imm6]. Register 31 reads as
// zr (SP/WSP are not allowed — use XZR/WZR). Shift — only lsl/lsr/asr;
// the 32-bit form limits the amount to 0..31 (see requireShift).
func (Builder) AddShift(rd, rn, rm Reg, imm Imm6, sh Shift) (Instr, error) {
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
		if err := requireClass(
			r.reg,
			"AddShift",
			r.op,
			"register 31 reads as zr — use XZR/WZR",
			classX,
			classW,
			classXZR,
			classWZR,
		); err != nil {
			return nil, err
		}
	}

	if err := requireWidth("AddShift", rd, rn, rm); err != nil {
		return nil, err
	}

	if err := requireShift(rd, "AddShift", imm, sh); err != nil {
		return nil, err
	}

	return newAddShift(base{}, rd.name(), rn.name(), rm.name(), imm.v, sh.String(), rd.Is64()), nil
}

func decodeAddShift(w uint32) Instr {
	return newAddShift(
		newBase(w),
		armRegName(w&0x1f, w>>31&1 == 1),
		armRegName(w>>5&0x1f, w>>31&1 == 1),
		armRegName(w>>16&0x1f, w>>31&1 == 1),
		w>>10&0x3f,
		shiftNames[w>>22&3],
		w>>31&1 == 1,
	)
}
