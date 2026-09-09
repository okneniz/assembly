package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Ldpsw — ldpsw rt, rt2, [rn, #imm7<<2] (offset form only).
type Ldpsw struct {
	base

	rt, rt2, rn string
	off         int64
}

// newLdpsw - the Ldpsw constructor: the struct is assembled only
// here (the Builder method and the decoder call it).
func newLdpsw(b base, rt string, rt2 string, rn string, off int64) Ldpsw {
	return Ldpsw{
		base: b,
		rt:   rt,
		rt2:  rt2,
		rn:   rn,
		off:  off,
	}
}

const ldpswEnc uint32 = 0x69400000

func (i Ldpsw) ObjDump(_ disasm.ViewCtx) string {
	if i.off == 0 {
		return fmt.Sprintf("ldpsw %s, %s, [%s]", i.rt, i.rt2, i.rn)
	}

	return fmt.Sprintf("ldpsw %s, %s, [%s, #%#x]", i.rt, i.rt2, i.rn, i.off)
}

func (i Ldpsw) Encode(w io.Writer) (int64, error) {
	rt, err := armRegNum(i.rt)
	if err != nil {
		return 0, fmt.Errorf("ldpsw: %w", err)
	}

	rt2, err := armRegNum(i.rt2)
	if err != nil {
		return 0, fmt.Errorf("ldpsw: %w", err)
	}

	rn, err := armRegNum(i.rn)
	if err != nil {
		return 0, fmt.Errorf("ldpsw: %w", err)
	}

	return writeWord(w, ldpswEnc|rt|rn<<5|rt2<<10|uint32(i.off>>2&0x7f)<<15)
}

// Ldpsw — ldpsw rt, rt2, [rn, #off]: sign-extending word load, rt
// and rt2 — x registers only (register 31 reads as zr), rn — x register
// or SP (register 31 in the base reads as sp); off — the signed imm7
// range scaled by 4.
func (Builder) Ldpsw(rt, rt2, rn Reg, off Off) (Instr, error) {
	if err := requireClass(rt, "Ldpsw", "rt", "x register (register 31 in rt reads as xzr)",
		classX, classXZR); err != nil {
		return nil, err
	}

	if err := requireClass(rt2, "Ldpsw", "rt2", "x register (register 31 in rt2 reads as xzr)",
		classX, classXZR); err != nil {
		return nil, err
	}

	if err := requireClass(
		rn,
		"Ldpsw",
		"rn",
		"x register or SP (register 31 in the base reads as sp)",
		classX,
		classSP,
	); err != nil {
		return nil, err
	}

	if err := requirePairOff("Ldpsw", off, 2); err != nil {
		return nil, err
	}

	return newLdpsw(base{}, rt.name(), rt2.name(), rn.name(), int64(off)), nil
}

func decodeLdpsw(w uint32) Instr {
	return newLdpsw(
		newBase(w),
		regNameX(w&0x1f),
		regNameX(w>>10&0x1f),
		regNameXSP(w>>5&0x1f),
		signExtendN(w>>15&0x7f, 7)<<2,
	)
}
