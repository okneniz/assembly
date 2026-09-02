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

const (
	SubShiftX uint32 = 0xCB000000
	SubShiftW uint32 = 0x4B000000
)

// SubShift — sub rd, rn, rm[, shift #imm6]. Register 31 reads as zr
// (SP/WSP are not allowed - use XZR/WZR). The shift is lsl/lsr/asr only;
// the 32-bit form limits the amount to 0..31 (see requireShift).
func (Builder) SubShift(rd, rn, rm Reg, imm Imm6, sh Shift) (Instr, error) {
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
			"SubShift",
			r.op,
			"register 31 reads as zr - use XZR/WZR",
			classX,
			classW,
			classXZR,
			classWZR,
		); err != nil {
			return nil, err
		}
	}

	if err := requireWidth("SubShift", rd, rn, rm); err != nil {
		return nil, err
	}

	if err := requireShift(rd, "SubShift", imm, sh); err != nil {
		return nil, err
	}

	return SubShift{
		rd:    rd.name(),
		rn:    rn.name(),
		rm:    rm.name(),
		imm6:  imm.v,
		shift: sh.String(),
		isf:   rd.Is64(),
	}, nil
}

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

func (i SubShift) Encode(w io.Writer, pc uint64) (int64, error) {
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

func (i SubShift) MarshalJSON() ([]byte, error) {
	return i.marshal("sub", i.ObjDump(disasm.DefaultViewCtx()), "Data processing - register",
		map[string]any{"Rd": i.rd, "Rn": i.rn, "Rm": i.rm, "imm6": i.imm6})
}

func decodeSubShift(w uint32, addr uint64) Instr {
	return SubShift{
		base:  newBase(addr, w),
		rd:    armRegName(w&0x1f, w>>31&1 == 1),
		rn:    armRegName(w>>5&0x1f, w>>31&1 == 1),
		rm:    armRegName(w>>16&0x1f, w>>31&1 == 1),
		imm6:  w >> 10 & 0x3f,
		shift: shiftNames[w>>22&3],
		isf:   w>>31&1 == 1,
	}
}
