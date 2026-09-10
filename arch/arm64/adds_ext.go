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

// newAddsExtBase - the AddsExt constructor for a ready embedded base (the
// string-operand layer): the struct is assembled only here.
func newAddsExtBase(b base, eb extBase) AddsExt {
	return AddsExt{
		base:    b,
		extBase: eb,
	}
}

// newAddsExt - the AddsExt constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newAddsExt(b base, rd Reg, rn Reg, rm Reg, ext string, imm3 uint32) (AddsExt, error) {
	err := requireClass(
		rd,
		"AddsExt",
		"rd",
		"register 31 reads as zr — use XZR/WZR (the cmn form)",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return AddsExt{}, err
	}

	err = requireClass(
		rn,
		"AddsExt",
		"rn",
		"register 31 reads as sp/wsp — use SP/WSP",
		classX,
		classW,
		classSP,
		classWSP,
	)

	if err != nil {
		return AddsExt{}, err
	}

	err = requireClass(
		rm,
		"AddsExt",
		"rm",
		"register 31 reads as sp/wsp — use SP/WSP",
		classX,
		classW,
		classSP,
		classWSP,
	)

	if err != nil {
		return AddsExt{}, err
	}

	err = requireWidth(
		"AddsExt",
		rd,
		rn,
		rm,
	)

	if err != nil {
		return AddsExt{}, err
	}

	if _, err := extNum(ext); err != nil {
		return AddsExt{}, fmt.Errorf("arm64.NewAddsExt: operand ext: %w", err)
	}

	if imm3 > 7 {
		return AddsExt{}, fmt.Errorf("arm64.NewAddsExt: operand imm3: %d is out of 0..7", imm3)
	}

	return AddsExt{
		base:    b,
		extBase: newExtBase(rd.bits(), rn.bits(), rm.bits(), ext, imm3, rd.Is64()),
	}, nil
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

func (Builder) AddsExt(rd, rn, rm Reg, ext string, imm3 uint32) (Instr, error) {
	return newAddsExt(base{}, rd, rn, rm, ext, imm3)
}

func decodeAddsExt(w uint32) (Instr, error) {
	in, err := newAddsExt(
		newBase(w),
		gprOf(w&0x1f, w>>31&1 == 1),
		gprOf(w>>5&0x1f, w>>31&1 == 1),
		gprOf(w>>16&0x1f, w>>31&1 == 1),
		extName(w>>13&7),
		w>>10&7,
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}
