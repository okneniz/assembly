package arm64

import (
	"errors"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AndShift — and rd, rn, rm[, shift #imm6].
type AndShift struct {
	base

	rd, rn, rm string
	imm6       uint32
	shift      string
	isf        bool
}

// newAndShift - the AndShift constructor: the struct is assembled only
// here (the Builder method and the decoder call it).
func newAndShift(
	b base,
	rd string,
	rn string,
	rm string,
	imm6 uint32,
	shift string,
	isf bool,
) AndShift {
	return AndShift{
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
	AndShiftX uint32 = 0x8A000000
	AndShiftW uint32 = 0x0A000000
)

func (i AndShift) ObjDump(_ disasm.ViewCtx) string {
	if i.imm6 != 0 {
		return fmt.Sprintf("and %s, %s, %s, %s #%d", i.rd, i.rn, i.rm, i.shift, i.imm6)
	}

	return fmt.Sprintf("and %s, %s, %s", i.rd, i.rn, i.rm)
}

func (i AndShift) Encode(w io.Writer) (int64, error) {
	match := AndShiftX
	if !i.isf {
		match = AndShiftW
	}

	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("and: %w", err)
	}

	sh, err := shiftNumByName(i.shift)
	if err != nil {
		return 0, fmt.Errorf("and: %w", err)
	}

	if i.imm6 > 63 {
		return 0, errors.New("and: imm6 out of range")
	}

	return writeWord(w, match|rd|rn<<5|i.imm6<<10|rm<<16|sh<<22)
}

// AndShift — and rd, rn, rm[, shift #imm6]. Register 31 reads as
// zr (SP/WSP are not allowed — use XZR/WZR); shift — lsl/lsr/asr/ror.
func (Builder) AndShift(rd, rn, rm Reg, imm Imm6, sh Shift) (Instr, error) {
	if err := requireClass(rd, "AndShift", "rd", "register 31 reads as zr — use XZR/WZR",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireClass(rn, "AndShift", "rn", "register 31 reads as zr — use XZR/WZR",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireClass(rm, "AndShift", "rm", "register 31 reads as zr — use XZR/WZR",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireWidth("AndShift", rd, rn, rm); err != nil {
		return nil, err
	}

	return newAndShift(base{}, rd.name(), rn.name(), rm.name(), imm.v, sh.String(), rd.Is64()), nil
}

func decodeAndShift(w uint32) Instr {
	return newAndShift(
		newBase(w),
		armRegName(w&0x1f, w>>31&1 == 1),
		armRegName(w>>5&0x1f, w>>31&1 == 1),
		armRegName(w>>16&0x1f, w>>31&1 == 1),
		w>>10&0x3f,
		shiftNames[w>>22&3],
		w>>31&1 == 1,
	)
}
