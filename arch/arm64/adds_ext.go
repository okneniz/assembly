package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AddsExt — adds rd, rn, rm{, ext #imm3}; pseudo: cmn (Rd=31).
type AddsExt struct {
	base
	extBase
}

// newAddsExt - the AddsExt constructor: the struct is assembled only
// here (the Builder method and the decoder call it).
func newAddsExt(b base, eb extBase) AddsExt {
	return AddsExt{
		base:    b,
		extBase: eb,
	}
}

const (
	AddsExtX uint32 = 0xAB200000
	AddsExtW uint32 = 0x2B200000
)

func (i AddsExt) ObjDump(_ disasm.ViewCtx) string {
	if i.rdNum == 31 {
		rnz := addSubRegName(i.rnNum, i.isf, false)
		return fmt.Sprintf(
			"cmn %s, %s%s",
			rnz,
			addSubRegName(i.rmNum, i.isf, false),
			i.extMod(true),
		)
	}

	return fmt.Sprintf("adds %s, %s, %s%s", addSubRegName(i.rdNum, i.isf, false),
		addSubRegName(i.rnNum, i.isf, false), addSubRegName(i.rmNum, i.isf, false), i.extMod(false))
}

func (i AddsExt) Encode(w io.Writer) (int64, error) {
	return i.extWrite(w, AddsExtX, AddsExtW, "adds")
}

// AddsExt — adds rd, rn, rm, ext #imm3 (cmn when Rd = zr). Rd:
// register 31 reads as zr; Rn/Rm — as sp/wsp; ext — uxtb..sxtx; imm3 — 0..7.
func (Builder) AddsExt(rd, rn, rm Reg, ext string, imm3 uint32) (Instr, error) {
	if err := requireClass(
		rd,
		"AddsExt",
		"rd",
		"register 31 reads as zr — use XZR/WZR (the cmn form)",
		classX,
		classW,
		classXZR,
		classWZR,
	); err != nil {
		return nil, err
	}

	if err := requireClass(rn, "AddsExt", "rn", "register 31 reads as sp/wsp — use SP/WSP",
		classX, classW, classSP, classWSP); err != nil {
		return nil, err
	}

	if err := requireClass(rm, "AddsExt", "rm", "register 31 reads as sp/wsp — use SP/WSP",
		classX, classW, classSP, classWSP); err != nil {
		return nil, err
	}

	if err := requireWidth("AddsExt", rd, rn, rm); err != nil {
		return nil, err
	}

	if _, err := extNum(ext); err != nil {
		return nil, fmt.Errorf("arm64.NewAddsExt: operand ext: %w", err)
	}

	if imm3 > 7 {
		return nil, fmt.Errorf("arm64.NewAddsExt: operand imm3: %d is out of 0..7", imm3)
	}

	return newAddsExt(
		base{},
		newExtBase(rd.bits(), rn.bits(), rm.bits(), ext, imm3, rd.Is64()),
	), nil
}

func decodeAddsExt(w uint32) Instr {
	return newAddsExt(newBase(w), decodeExtBase(w))
}
