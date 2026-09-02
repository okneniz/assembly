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

const (
	AddsShiftX uint32 = 0xAB000000
	AddsShiftW uint32 = 0x2B000000
)

// AddsShift — adds rd, rn, rm[, shift #imm6] (cmn when Rd = zr).
// Register 31 reads as zr (SP/WSP are not allowed — use XZR/WZR). Shift —
// only lsl/lsr/asr; the 32-bit form limits the amount to 0..31 (see requireShift).
func (Builder) AddsShift(rd, rn, rm Reg, imm Imm6, sh Shift) (Instr, error) {
	if err := requireClass(rd, "AddsShift", "rd", "register 31 reads as zr — use XZR/WZR",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireClass(rn, "AddsShift", "rn", "register 31 reads as zr — use XZR/WZR",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireClass(rm, "AddsShift", "rm", "register 31 reads as zr — use XZR/WZR",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireWidth("AddsShift", rd, rn, rm); err != nil {
		return nil, err
	}

	if err := requireShift(rd, "AddsShift", imm, sh); err != nil {
		return nil, err
	}

	return AddsShift{
		rd:    rd.name(),
		rn:    rn.name(),
		rm:    rm.name(),
		imm6:  imm.v,
		shift: sh.String(),
		isf:   rd.Is64(),
	}, nil
}

func decodeAddsShift(w uint32, addr uint64) Instr {
	return AddsShift{
		base:  newBase(addr, w),
		rd:    armRegName(w&0x1f, w>>31&1 == 1),
		rn:    armRegName(w>>5&0x1f, w>>31&1 == 1),
		rm:    armRegName(w>>16&0x1f, w>>31&1 == 1),
		imm6:  w >> 10 & 0x3f,
		shift: shiftNames[w>>22&3],
		isf:   w>>31&1 == 1,
	}
}

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

func (i AddsShift) Encode(w io.Writer, pc uint64) (int64, error) {
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

func (i AddsShift) MarshalJSON() ([]byte, error) {
	return i.marshal("adds", i.ObjDump(disasm.DefaultViewCtx()), "Data processing - register",
		map[string]any{"Rd": i.rd, "Rn": i.rn, "Rm": i.rm, "imm6": i.imm6})
}
