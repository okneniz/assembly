package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Ldrsb — ldrsb rt, [rn, #imm12] (Rt — X; the W form is not decoded).
type Ldrsb struct {
	base

	rt, rn string
	off    int64
}

const ldrsbEnc uint32 = 0x39800000

// Ldrsb — ldrsb rt, [rn, #off]: sign-extending byte load, rt — x
// register only (register 31 reads as xzr), rn — x register or SP
// (register 31 in the base reads as sp); the offset is an unscaled
// imm12 (0..0xfff).
func (Builder) Ldrsb(rt, rn Reg, off Off) (Instr, error) {
	if err := requireClass(rt, "Ldrsb", "rt", "x register (register 31 in rt reads as xzr)",
		classX, classXZR); err != nil {
		return nil, err
	}

	if err := requireClass(
		rn,
		"Ldrsb",
		"rn",
		"x register or SP (register 31 in the base reads as sp)",
		classX,
		classSP,
	); err != nil {
		return nil, err
	}

	if err := requireOff("Ldrsb", off, 0); err != nil {
		return nil, err
	}

	return Ldrsb{
		rt:  rt.name(),
		rn:  rn.name(),
		off: int64(off),
	}, nil
}

func decodeLdrsb(w uint32, addr uint64) Instr {
	return Ldrsb{
		base: newBase(addr, w),
		rt:   regNameX(w & 0x1f),
		rn:   regNameXSP(w >> 5 & 0x1f),
		off:  int64(w >> 10 & 0xfff),
	}
}

func (i Ldrsb) ObjDump(_ disasm.ViewCtx) string {
	if i.off == 0 {
		return fmt.Sprintf("ldrsb %s, [%s]", i.rt, i.rn)
	}

	return fmt.Sprintf("ldrsb %s, [%s, #0x%x]", i.rt, i.rn, i.off)
}

func (i Ldrsb) Encode(w io.Writer, pc uint64) (int64, error) {
	return lsSignedWrite(w, ldrsbEnc, i.rt, i.rn, i.off, "ldrsb")
}

func (i Ldrsb) MarshalJSON() ([]byte, error) {
	return i.marshal(
		"ldrsb",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Load/Store",
		map[string]any{"Rt": i.rt, "Rn": i.rn},
	)
}
