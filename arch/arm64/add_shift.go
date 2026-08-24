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

const (
	AddShiftX uint32 = 0x8B000000
	AddShiftW uint32 = 0x0B000000
)

// NewAddShift — add rd, rn, rm[, shift #imm6]. Register 31 reads as
// zr (SP/WSP are not allowed — use XZR/WZR). Shift — only lsl/lsr/asr;
// the 32-bit form limits the amount to 0..31 (see requireShift).
func NewAddShift(rd, rn, rm Reg, imm Imm6, sh Shift) (Instr, error) {
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

	return AddShift{
		rd:    rd.name(),
		rn:    rn.name(),
		rm:    rm.name(),
		imm6:  imm.v,
		shift: sh.String(),
		isf:   rd.Is64(),
	}, nil
}

func (i AddShift) ObjDump(_ disasm.ViewCtx) string {
	if i.imm6 != 0 {
		return fmt.Sprintf("add %s, %s, %s, %s #%d", i.rd, i.rn, i.rm, i.shift, i.imm6)
	}

	return fmt.Sprintf("add %s, %s, %s", i.rd, i.rn, i.rm)
}

func (i AddShift) Encode(w io.Writer, pc uint64) (int64, error) {
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

func (i AddShift) MarshalJSON() ([]byte, error) {
	return i.marshal("add", i.ObjDump(disasm.DefaultViewCtx()), "Data processing - register",
		map[string]any{"Rd": i.rd, "Rn": i.rn, "Rm": i.rm, "imm6": i.imm6})
}

func decodeAddShift(w uint32, addr uint64) Instr {
	return AddShift{
		base:  newBase(addr, w),
		rd:    armRegName(w&0x1f, w>>31&1 == 1),
		rn:    armRegName(w>>5&0x1f, w>>31&1 == 1),
		rm:    armRegName(w>>16&0x1f, w>>31&1 == 1),
		imm6:  w >> 10 & 0x3f,
		shift: shiftNames[w>>22&3],
		isf:   w>>31&1 == 1,
	}
}
