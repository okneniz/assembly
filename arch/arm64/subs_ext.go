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

const (
	SubsExtX uint32 = 0xEB200000
	SubsExtW uint32 = 0x6B200000
)

// SubsExt — subs rd, rn, rm, ext #imm3 (cmp when Rd = zr). Rd:
// register 31 reads as zr; Rn/Rm — as sp/wsp; ext — uxtb..sxtx;
// imm3 — 0..7.
func (Builder) SubsExt(rd, rn, rm Reg, ext string, imm3 uint32) (Instr, error) {
	if err := requireClass(
		rd,
		"SubsExt",
		"rd",
		"register 31 reads as zr — use XZR/WZR (the cmp form)",
		classX,
		classW,
		classXZR,
		classWZR,
	); err != nil {
		return nil, err
	}

	if err := requireClass(rn, "SubsExt", "rn", "register 31 reads as sp/wsp — use SP/WSP",
		classX, classW, classSP, classWSP); err != nil {
		return nil, err
	}

	if err := requireClass(rm, "SubsExt", "rm", "register 31 reads as sp/wsp — use SP/WSP",
		classX, classW, classSP, classWSP); err != nil {
		return nil, err
	}

	if err := requireWidth("SubsExt", rd, rn, rm); err != nil {
		return nil, err
	}

	if _, err := extNum(ext); err != nil {
		return nil, fmt.Errorf("arm64.NewSubsExt: operand ext: %w", err)
	}

	if imm3 > 7 {
		return nil, fmt.Errorf("arm64.NewSubsExt: operand imm3: %d is out of 0..7", imm3)
	}

	return SubsExt{extBase: newExtBase(rd.bits(), rn.bits(), rm.bits(), ext, imm3, rd.Is64())}, nil
}

func decodeSubsExt(w uint32, addr uint64) Instr {
	return SubsExt{
		base:    newBase(addr, w),
		extBase: decodeExtBase(w),
	}
}

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

func (i SubsExt) Encode(w io.Writer, pc uint64) (int64, error) {
	return i.extWrite(w, SubsExtX, SubsExtW, "subs")
}

func (i SubsExt) MarshalJSON() ([]byte, error) {
	return i.marshal(
		"subs",
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
