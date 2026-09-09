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

// newEorShift - the EorShift constructor: the struct is assembled only
// here (the Builder method and the decoder call it).
func newEorShift(
	b base,
	rd string,
	rn string,
	rm string,
	imm6 uint32,
	shift string,
	isf bool,
) EorShift {
	return EorShift{
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

// EorShift — eor rd, rn, rm[, shift #imm6]. Register 31 reads as
// zr (SP/WSP are not allowed — use XZR/WZR); shift — lsl/lsr/asr/ror.
func (Builder) EorShift(rd, rn, rm Reg, imm Imm6, sh Shift) (Instr, error) {
	if err := requireClass(rd, "EorShift", "rd", "register 31 reads as zr — use XZR/WZR",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireClass(rn, "EorShift", "rn", "register 31 reads as zr — use XZR/WZR",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireClass(rm, "EorShift", "rm", "register 31 reads as zr — use XZR/WZR",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireWidth("EorShift", rd, rn, rm); err != nil {
		return nil, err
	}

	return newEorShift(base{}, rd.name(), rn.name(), rm.name(), imm.v, sh.String(), rd.Is64()), nil
}

func decodeEorShift(w uint32) Instr {
	return newEorShift(
		newBase(w),
		armRegName(w&0x1f, w>>31&1 == 1),
		armRegName(w>>5&0x1f, w>>31&1 == 1),
		armRegName(w>>16&0x1f, w>>31&1 == 1),
		w>>10&0x3f,
		shiftNames[w>>22&3],
		w>>31&1 == 1,
	)
}
