package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Ccmp — ccmp rn, rm, #imm, cond.
type Ccmp struct {
	base

	rn, rm string
	immVal uint32
	cond   string
}

// newCcmp - the Ccmp constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newCcmp(b base, rn Reg, rm Reg, nzcv uint32, cond string) (Ccmp, error) {
	err := requireClass(
		rn,
		"Ccmp",
		"rn",
		"register 31 reads as zr — use XZR (only the 64-bit form)",
		classX,
		classXZR,
	)

	if err != nil {
		return Ccmp{}, err
	}

	err = requireClass(
		rm,
		"Ccmp",
		"rm",
		"register 31 reads as zr — use XZR (only the 64-bit form)",
		classX,
		classXZR,
	)

	if err != nil {
		return Ccmp{}, err
	}

	if nzcv > 0xf {
		return Ccmp{}, fmt.Errorf("arm64.NewCcmp: operand nzcv: %#x is out of 0..0xf", nzcv)
	}

	if _, err := condNum(cond); err != nil {
		return Ccmp{}, fmt.Errorf("arm64.NewCcmp: operand cond: %w", err)
	}

	return Ccmp{
		base:   b,
		rn:     rn.name(),
		rm:     rm.name(),
		immVal: nzcv,
		cond:   cond,
	}, nil
}

const ccmpX uint32 = 0xFA400000

func (i Ccmp) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ccmp %s, %s, #0x%x, %s", i.rn, i.rm, i.immVal, i.cond)
}

func (i Ccmp) Encode(w io.Writer) (int64, error) {
	rn, rm, err := regNums2(i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("ccmp: %w", err)
	}

	c, err := condNum(i.cond)
	if err != nil {
		return 0, fmt.Errorf("ccmp: %w", err)
	}

	if i.immVal > 0xf {
		return 0, fmt.Errorf("ccmp: imm %#x out of range", i.immVal)
	}

	return writeWord(w, ccmpX|i.immVal|rn<<5|c<<12|rm<<16)
}

func (Builder) Ccmp(rn, rm Reg, nzcv uint32, cond string) (Instr, error) {
	return newCcmp(base{}, rn, rm, nzcv, cond)
}

func decodeCcmp(w uint32) (Instr, error) {
	in, err := newCcmp(
		newBase(w),
		gprOf(w>>5&0x1f, true),
		gprOf(w>>16&0x1f, true),
		w&0xf,
		condName(w>>12&0xf),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}
