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

// newCsel - the Csel constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newCsel(b base, rd Reg, rn Reg, rm Reg, cond string) (Csel, error) {
	err := requireClass(
		rd,
		"Csel",
		"rd",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return Csel{}, err
	}

	err = requireClass(
		rn,
		"Csel",
		"rn",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return Csel{}, err
	}

	err = requireClass(
		rm,
		"Csel",
		"rm",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return Csel{}, err
	}

	err = requireWidth(
		"Csel",
		rd,
		rn,
		rm,
	)

	if err != nil {
		return Csel{}, err
	}

	if _, err := condNum(cond); err != nil {
		return Csel{}, fmt.Errorf("arm64.NewCsel: operand cond: %w", err)
	}

	return Csel{
		base: b,
		rd:   rd.name(),
		rn:   rn.name(),
		rm:   rm.name(),
		cond: cond,
	}, nil
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

func (Builder) Csel(rd, rn, rm Reg, cond string) (Instr, error) {
	return newCsel(base{}, rd, rn, rm, cond)
}

func decodeCsel(w uint32) (Instr, error) {
	in, err := newCsel(
		newBase(w),
		gprOf(w&0x1f, w>>31&1 == 1),
		gprOf(w>>5&0x1f, w>>31&1 == 1),
		gprOf(w>>16&0x1f, w>>31&1 == 1),
		condName(w>>12&0xf),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}
