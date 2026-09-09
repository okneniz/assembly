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

// newCcmp - the Ccmp constructor: the struct is assembled only
// here (the Builder method and the decoder call it).
func newCcmp(b base, rn string, rm string, immVal uint32, cond string) Ccmp {
	return Ccmp{
		base:   b,
		rn:     rn,
		rm:     rm,
		immVal: immVal,
		cond:   cond,
	}
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

// Ccmp — ccmp rn, rm, #nzcv, cond. Only the 64-bit form;
// register 31 reads as zr (use XZR); nzcv — 0..0xf; cond — the standard
// condition names (eq/ne/.../nv, see condNum).
func (Builder) Ccmp(rn, rm Reg, nzcv uint32, cond string) (Instr, error) {
	if err := requireClass(rn, "Ccmp", "rn",
		"register 31 reads as zr — use XZR (only the 64-bit form)", classX, classXZR); err != nil {
		return nil, err
	}

	if err := requireClass(rm, "Ccmp", "rm",
		"register 31 reads as zr — use XZR (only the 64-bit form)", classX, classXZR); err != nil {
		return nil, err
	}

	if nzcv > 0xf {
		return nil, fmt.Errorf("arm64.NewCcmp: operand nzcv: %#x is out of 0..0xf", nzcv)
	}

	if _, err := condNum(cond); err != nil {
		return nil, fmt.Errorf("arm64.NewCcmp: operand cond: %w", err)
	}

	return newCcmp(base{}, rn.name(), rm.name(), nzcv, cond), nil
}

func decodeCcmp(w uint32) Instr {
	return newCcmp(
		newBase(w),
		armRegName(w>>5&0x1f, true),
		armRegName(w>>16&0x1f, true),
		w&0xf,
		condName(w>>12&0xf),
	)
}
