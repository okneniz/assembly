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

const ldpswEnc uint32 = 0x69400000

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

	return Ldpsw{
		rt:  rt.name(),
		rt2: rt2.name(),
		rn:  rn.name(),
		off: int64(off),
	}, nil
}

func decodeLdpsw(w uint32, addr uint64) Instr {
	return Ldpsw{
		base: newBase(addr, w),
		rt:   regNameX(w & 0x1f),
		rt2:  regNameX(w >> 10 & 0x1f),
		rn:   regNameXSP(w >> 5 & 0x1f),
		off:  signExtendN(w>>15&0x7f, 7) << 2,
	}
}

func (i Ldpsw) ObjDump(_ disasm.ViewCtx) string {
	if i.off == 0 {
		return fmt.Sprintf("ldpsw %s, %s, [%s]", i.rt, i.rt2, i.rn)
	}

	return fmt.Sprintf("ldpsw %s, %s, [%s, #%#x]", i.rt, i.rt2, i.rn, i.off)
}

func (i Ldpsw) Encode(w io.Writer, pc uint64) (int64, error) {
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

func (i Ldpsw) MarshalJSON() ([]byte, error) {
	return i.marshal("ldpsw", i.ObjDump(disasm.DefaultViewCtx()), "Load/Store",
		map[string]any{"Rt": i.rt, "Rt2": i.rt2, "Rn": i.rn})
}
