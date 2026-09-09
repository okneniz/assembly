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

// newSubExt - the SubExt constructor: the struct is assembled only
// here (the Builder method and the decoder call it).
func newSubExt(b base, eb extBase) SubExt {
	return SubExt{
		base:    b,
		extBase: eb,
	}
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

// SubExt — sub rd, rn, rm, ext #imm3. Register 31 reads as
// sp/wsp; ext — uxtb/uxth/uxtw/uxtx/sxtb/sxth/sxtw/sxtx; imm3 — 0..7.
func (Builder) SubExt(rd, rn, rm Reg, ext string, imm3 uint32) (Instr, error) {
	if err := requireClass(rd, "SubExt", "rd", "register 31 reads as sp/wsp — use SP/WSP",
		classX, classW, classSP, classWSP); err != nil {
		return nil, err
	}

	if err := requireClass(rn, "SubExt", "rn", "register 31 reads as sp/wsp — use SP/WSP",
		classX, classW, classSP, classWSP); err != nil {
		return nil, err
	}

	if err := requireClass(rm, "SubExt", "rm", "register 31 reads as sp/wsp — use SP/WSP",
		classX, classW, classSP, classWSP); err != nil {
		return nil, err
	}

	if err := requireWidth("SubExt", rd, rn, rm); err != nil {
		return nil, err
	}

	if _, err := extNum(ext); err != nil {
		return nil, fmt.Errorf("arm64.NewSubExt: operand ext: %w", err)
	}

	if imm3 > 7 {
		return nil, fmt.Errorf("arm64.NewSubExt: operand imm3: %d is out of 0..7", imm3)
	}

	return newSubExt(base{}, newExtBase(rd.bits(), rn.bits(), rm.bits(), ext, imm3, rd.Is64())), nil
}

func decodeSubExt(w uint32) Instr {
	return newSubExt(newBase(w), decodeExtBase(w))
}
