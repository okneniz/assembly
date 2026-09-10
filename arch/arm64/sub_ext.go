package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// SubExt — sub rd, rn, rm{, ext #imm3}.
type SubExt struct {
	base
	extBase
}

// newSubExtBase - the SubExt constructor for a ready embedded base (the
// string-operand layer): the struct is assembled only here.
func newSubExtBase(b base, eb extBase) SubExt {
	return SubExt{
		base:    b,
		extBase: eb,
	}
}

// newSubExt - the SubExt constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newSubExt(b base, rd Reg, rn Reg, rm Reg, ext string, imm3 uint32) (SubExt, error) {
	err := requireClass(
		rd,
		"SubExt",
		"rd",
		"register 31 reads as sp/wsp — use SP/WSP",
		classX,
		classW,
		classSP,
		classWSP,
	)

	if err != nil {
		return SubExt{}, err
	}

	err = requireClass(
		rn,
		"SubExt",
		"rn",
		"register 31 reads as sp/wsp — use SP/WSP",
		classX,
		classW,
		classSP,
		classWSP,
	)

	if err != nil {
		return SubExt{}, err
	}

	err = requireClass(
		rm,
		"SubExt",
		"rm",
		"register 31 reads as sp/wsp — use SP/WSP",
		classX,
		classW,
		classSP,
		classWSP,
	)

	if err != nil {
		return SubExt{}, err
	}

	err = requireWidth(
		"SubExt",
		rd,
		rn,
		rm,
	)

	if err != nil {
		return SubExt{}, err
	}

	if _, err := extNum(ext); err != nil {
		return SubExt{}, fmt.Errorf("arm64.NewSubExt: operand ext: %w", err)
	}

	if imm3 > 7 {
		return SubExt{}, fmt.Errorf("arm64.NewSubExt: operand imm3: %d is out of 0..7", imm3)
	}

	return SubExt{
		base:    b,
		extBase: newExtBase(rd.bits(), rn.bits(), rm.bits(), ext, imm3, rd.Is64()),
	}, nil
}

const (
	SubExtX uint32 = 0xCB200000
	SubExtW uint32 = 0x4B200000
)

func (i SubExt) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("sub %s, %s, %s%s", addSubRegName(i.rdNum, i.isf, false),
		addSubRegName(i.rnNum, i.isf, false), addSubRegName(i.rmNum, i.isf, false), i.extMod(false))
}

func (i SubExt) Encode(w io.Writer) (int64, error) {
	return i.extWrite(w, SubExtX, SubExtW, "sub")
}

func (Builder) SubExt(rd, rn, rm Reg, ext string, imm3 uint32) (Instr, error) {
	return newSubExt(base{}, rd, rn, rm, ext, imm3)
}

func decodeSubExt(w uint32) (Instr, error) {
	in, err := newSubExt(
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
