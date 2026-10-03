package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// SubExt — sub rd, rn, rm{, ext #imm3}.
type SubExt struct {
	extBase
}

// newSubExtBase - the SubExt constructor for a ready embedded base (the
// string-operand layer): the struct is assembled only here.
func newSubExtBase(eb extBase) SubExt {
	return SubExt{
		extBase: eb,
	}
}

// newSubExt - the SubExt constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newSubExt(rd Reg, rn Reg, rm Reg, ext string, imm3 uint32) (SubExt, error) {
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
		"register 31 reads as zr — use XZR/WZR (sp is not an operand here)",
		classX,
		classW,
		classXZR,
		classWZR,
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

	if err := requireExtWidth("SubExt", rd.Is64(), ext); err != nil {
		return SubExt{}, err
	}

	if imm3 > 7 {
		return SubExt{}, fmt.Errorf("arm64.NewSubExt: operand imm3: %d is out of 0..7", imm3)
	}

	return SubExt{
		extBase: newExtBase(rd.bits(), rn.bits(), rm.bits(), ext, imm3, rd.Is64()),
	}, nil
}

const (
	SubExtX uint32 = 0xCB200000
	SubExtW uint32 = 0x4B200000
)

func (i SubExt) Encode(w io.Writer) (int64, error) {
	return i.extWrite(w, SubExtX, SubExtW, "sub")
}

func (i SubExt) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("sub %s, %s, %s%s", addSubRegName(i.rdNum, i.isf, false),
		addSubRegName(i.rnNum, i.isf, false), addSubRegName(i.rmNum, i.isf, false), i.extMod(false))
}
