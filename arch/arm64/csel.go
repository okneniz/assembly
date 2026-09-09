package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Csel — csel rd, rn, rm, cond. Csinc/Csinv/Csneg — the same encodings
// with the inverse aliases cset/csetm/cinc/cinv/cneg.
type Csel struct {
	base

	rd, rn, rm string
	cond       string
}

// newCsel - the Csel constructor: the struct is assembled only
// here (the Builder method and the decoder call it).
func newCsel(b base, rd string, rn string, rm string, cond string) Csel {
	return Csel{
		base: b,
		rd:   rd,
		rn:   rn,
		rm:   rm,
		cond: cond,
	}
}

const (
	cselX uint32 = 0x9A800000
	cselW uint32 = 0x1A800000
)

func (i Csel) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("csel %s, %s, %s, %s", i.rd, i.rn, i.rm, i.cond)
}

func (i Csel) Encode(w io.Writer) (int64, error) {
	return cselWrite(w, i, cselX, cselW, "csel")
}

// cselWrite — the common encoding skeleton of the csel family.
func cselWrite(w io.Writer, i Csel, matchX, matchW uint32, name string) (int64, error) {
	match, err := sfMatch(i.rd, matchX, matchW)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", name, err)
	}

	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", name, err)
	}

	c, err := condNum(i.cond)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", name, err)
	}

	return writeWord(w, match|rd|rn<<5|c<<12|rm<<16)
}

// Csel — csel rd, rn, rm, cond. Register 31 reads as zr (SP/WSP
// are not allowed — use XZR/WZR); the width is shared by all three
// registers; cond — the standard condition names (eq/ne/.../nv, see
// condNum).
func (Builder) Csel(rd, rn, rm Reg, cond string) (Instr, error) {
	if err := requireClass(rd, "Csel", "rd", "register 31 reads as zr — use XZR/WZR",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireClass(rn, "Csel", "rn", "register 31 reads as zr — use XZR/WZR",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireClass(rm, "Csel", "rm", "register 31 reads as zr — use XZR/WZR",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireWidth("Csel", rd, rn, rm); err != nil {
		return nil, err
	}

	if _, err := condNum(cond); err != nil {
		return nil, fmt.Errorf("arm64.NewCsel: operand cond: %w", err)
	}

	return newCsel(base{}, rd.name(), rn.name(), rm.name(), cond), nil
}

func decodeCsel(w uint32) Instr {
	return newCsel(
		newBase(w),
		armRegName(w&0x1f, w>>31&1 == 1),
		armRegName(w>>5&0x1f, w>>31&1 == 1),
		armRegName(w>>16&0x1f, w>>31&1 == 1),
		condName(w>>12&0xf),
	)
}
