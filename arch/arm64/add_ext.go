package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AddExt — add rd, rn, rm{, ext #imm3}.
type AddExt struct {
	base
	extBase
}

// newAddExtBase - the AddExt constructor for a ready embedded base (the
// string-operand layer): the struct is assembled only here.
func newAddExtBase(b base, eb extBase) AddExt {
	return AddExt{
		base:    b,
		extBase: eb,
	}
}

// newAddExt - the AddExt constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newAddExt(b base, rd Reg, rn Reg, rm Reg, ext string, imm3 uint32) (AddExt, error) {
	err := requireClass(
		rd,
		"AddExt",
		"rd",
		"register 31 reads as sp/wsp — use SP/WSP",
		classX,
		classW,
		classSP,
		classWSP,
	)

	if err != nil {
		return AddExt{}, err
	}

	err = requireClass(
		rn,
		"AddExt",
		"rn",
		"register 31 reads as sp/wsp — use SP/WSP",
		classX,
		classW,
		classSP,
		classWSP,
	)

	if err != nil {
		return AddExt{}, err
	}

	err = requireClass(
		rm,
		"AddExt",
		"rm",
		"register 31 reads as sp/wsp — use SP/WSP",
		classX,
		classW,
		classSP,
		classWSP,
	)

	if err != nil {
		return AddExt{}, err
	}

	err = requireWidth(
		"AddExt",
		rd,
		rn,
		rm,
	)

	if err != nil {
		return AddExt{}, err
	}

	if _, err := extNum(ext); err != nil {
		return AddExt{}, fmt.Errorf("arm64.NewAddExt: operand ext: %w", err)
	}

	if imm3 > 7 {
		return AddExt{}, fmt.Errorf("arm64.NewAddExt: operand imm3: %d is out of 0..7", imm3)
	}

	return AddExt{
		base:    b,
		extBase: newExtBase(rd.bits(), rn.bits(), rm.bits(), ext, imm3, rd.Is64()),
	}, nil
}

const (
	AddExtX uint32 = 0x8B200000
	AddExtW uint32 = 0x0B200000
)

func (i AddExt) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("add %s, %s, %s%s", addSubRegName(i.rdNum, i.isf, false),
		addSubRegName(i.rnNum, i.isf, false), addSubRegName(i.rmNum, i.isf, false), i.extMod(false))
}

func (i AddExt) Encode(w io.Writer) (int64, error) {
	return i.extWrite(w, AddExtX, AddExtW, "add")
}

func (Builder) AddExt(rd, rn, rm Reg, ext string, imm3 uint32) (Instr, error) {
	return newAddExt(base{}, rd, rn, rm, ext, imm3)
}

func decodeAddExt(w uint32) (Instr, error) {
	in, err := newAddExt(
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
