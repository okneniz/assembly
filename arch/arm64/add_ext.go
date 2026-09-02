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

const (
	AddExtX uint32 = 0x8B200000
	AddExtW uint32 = 0x0B200000
)

// AddExt — add rd, rn, rm, ext #imm3. Register 31 reads as
// sp/wsp; ext — uxtb/uxth/uxtw/uxtx/sxtb/sxth/sxtw/sxtx; imm3 — 0..7.
func (Builder) AddExt(rd, rn, rm Reg, ext string, imm3 uint32) (Instr, error) {
	if err := requireClass(rd, "AddExt", "rd", "register 31 reads as sp/wsp — use SP/WSP",
		classX, classW, classSP, classWSP); err != nil {
		return nil, err
	}

	if err := requireClass(rn, "AddExt", "rn", "register 31 reads as sp/wsp — use SP/WSP",
		classX, classW, classSP, classWSP); err != nil {
		return nil, err
	}

	if err := requireClass(rm, "AddExt", "rm", "register 31 reads as sp/wsp — use SP/WSP",
		classX, classW, classSP, classWSP); err != nil {
		return nil, err
	}

	if err := requireWidth("AddExt", rd, rn, rm); err != nil {
		return nil, err
	}

	if _, err := extNum(ext); err != nil {
		return nil, fmt.Errorf("arm64.NewAddExt: operand ext: %w", err)
	}

	if imm3 > 7 {
		return nil, fmt.Errorf("arm64.NewAddExt: operand imm3: %d is out of 0..7", imm3)
	}

	return AddExt{extBase: newExtBase(rd.bits(), rn.bits(), rm.bits(), ext, imm3, rd.Is64())}, nil
}

func decodeAddExt(w uint32, addr uint64) Instr {
	return AddExt{
		base:    newBase(addr, w),
		extBase: decodeExtBase(w),
	}
}

func (i AddExt) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("add %s, %s, %s%s", addSubRegName(i.rdNum, i.isf, false),
		addSubRegName(i.rnNum, i.isf, false), addSubRegName(i.rmNum, i.isf, false), i.extMod(false))
}

func (i AddExt) Encode(w io.Writer, pc uint64) (int64, error) {
	return i.extWrite(w, AddExtX, AddExtW, "add")
}

func (i AddExt) MarshalJSON() ([]byte, error) {
	return i.marshal(
		"add",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Data processing - register",
		map[string]any{
			"Rd":     i.rdNum,
			"Rn":     i.rnNum,
			"Rm":     i.rmNum,
			"option": i.option,
			"imm3":   i.imm3,
		},
	)
}
