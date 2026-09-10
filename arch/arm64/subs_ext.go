package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// SubsExt — subs rd, rn, rm{, ext #imm3}; pseudo: cmp (Rd=31).
type SubsExt struct {
	base
	extBase
}

// newSubsExtBase - the SubsExt constructor for a ready embedded base (the
// string-operand layer): the struct is assembled only here.
func newSubsExtBase(b base, eb extBase) SubsExt {
	return SubsExt{
		base:    b,
		extBase: eb,
	}
}

// newSubsExt - the SubsExt constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newSubsExt(b base, rd Reg, rn Reg, rm Reg, ext string, imm3 uint32) (SubsExt, error) {
	err := requireClass(
		rd,
		"SubsExt",
		"rd",
		"register 31 reads as zr — use XZR/WZR (the cmp form)",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return SubsExt{}, err
	}

	err = requireClass(
		rn,
		"SubsExt",
		"rn",
		"register 31 reads as sp/wsp — use SP/WSP",
		classX,
		classW,
		classSP,
		classWSP,
	)

	if err != nil {
		return SubsExt{}, err
	}

	err = requireClass(
		rm,
		"SubsExt",
		"rm",
		"register 31 reads as sp/wsp — use SP/WSP",
		classX,
		classW,
		classSP,
		classWSP,
	)

	if err != nil {
		return SubsExt{}, err
	}

	err = requireWidth(
		"SubsExt",
		rd,
		rn,
		rm,
	)

	if err != nil {
		return SubsExt{}, err
	}

	if _, err := extNum(ext); err != nil {
		return SubsExt{}, fmt.Errorf("arm64.NewSubsExt: operand ext: %w", err)
	}

	if imm3 > 7 {
		return SubsExt{}, fmt.Errorf("arm64.NewSubsExt: operand imm3: %d is out of 0..7", imm3)
	}

	return SubsExt{
		base:    b,
		extBase: newExtBase(rd.bits(), rn.bits(), rm.bits(), ext, imm3, rd.Is64()),
	}, nil
}

const (
	SubsExtX uint32 = 0xEB200000
	SubsExtW uint32 = 0x6B200000
)

func (i SubsExt) ObjDump(_ disasm.ViewCtx) string {
	if i.rdNum == 31 {
		rnz := addSubRegName(i.rnNum, i.isf, false)
		return fmt.Sprintf(
			"cmp %s, %s%s",
			rnz,
			addSubRegName(i.rmNum, i.isf, false),
			i.extMod(true),
		)
	}

	return fmt.Sprintf("subs %s, %s, %s%s", addSubRegName(i.rdNum, i.isf, false),
		addSubRegName(i.rnNum, i.isf, false), addSubRegName(i.rmNum, i.isf, false), i.extMod(false))
}

func (i SubsExt) Encode(w io.Writer) (int64, error) {
	return i.extWrite(w, SubsExtX, SubsExtW, "subs")
}

func (Builder) SubsExt(rd, rn, rm Reg, ext string, imm3 uint32) (Instr, error) {
	return newSubsExt(base{}, rd, rn, rm, ext, imm3)
}

func decodeSubsExt(w uint32) (Instr, error) {
	in, err := newSubsExt(
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
