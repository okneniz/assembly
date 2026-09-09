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

// newAddsShift - the AddsShift constructor: the struct is assembled only
// here (the Builder method and the decoder call it).
func newAddsShift(
	b base,
	rd string,
	rn string,
	rm string,
	imm6 uint32,
	shift string,
	isf bool,
) AddsShift {
	return AddsShift{
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
	AddsShiftX uint32 = 0xAB000000
	AddsShiftW uint32 = 0x2B000000
)

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

func (i AddsShift) Encode(w io.Writer) (int64, error) {
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

	return newAddsShift(base{}, rd.name(), rn.name(), rm.name(), imm.v, sh.String(), rd.Is64()), nil
}

func decodeAddsShift(w uint32) Instr {
	return newAddsShift(
		newBase(w),
		armRegName(w&0x1f, w>>31&1 == 1),
		armRegName(w>>5&0x1f, w>>31&1 == 1),
		armRegName(w>>16&0x1f, w>>31&1 == 1),
		w>>10&0x3f,
		shiftNames[w>>22&3],
		w>>31&1 == 1,
	)
}
