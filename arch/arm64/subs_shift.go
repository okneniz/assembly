package arm64

import (
	"errors"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// SubsShift - subs rd, rn, rm[, shift #imm6]; pseudo: cmp (Rd = zr), negs (Rn = zr).
type SubsShift struct {
	base

	rd, rn, rm string
	imm6       uint32
	shift      string // lsl/lsr/asr/ror
	isf        bool
}

const (
	SubsShiftX uint32 = 0xEB000000
	SubsShiftW uint32 = 0x6B000000
)

func decodeSubsShift(w uint32, addr uint64) Instr {
	return SubsShift{
		base:  newBase(addr, w),
		rd:    armRegName(w&0x1f, w>>31&1 == 1),
		rn:    armRegName(w>>5&0x1f, w>>31&1 == 1),
		rm:    armRegName(w>>16&0x1f, w>>31&1 == 1),
		imm6:  w >> 10 & 0x3f,
		shift: shiftNames[w>>22&3],
		isf:   w>>31&1 == 1,
	}
}

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

func (i SubsShift) Encode(w io.Writer, pc uint64) (int64, error) {
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

func (i SubsShift) MarshalJSON() ([]byte, error) {
	return i.marshal("subs", i.ObjDump(disasm.DefaultViewCtx()), "Data processing - register",
		map[string]any{"Rd": i.rd, "Rn": i.rn, "Rm": i.rm, "imm6": i.imm6})
}
