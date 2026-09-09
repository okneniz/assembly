package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Ldrsh — ldrsh rt, [rn, #imm12<<1].
type Ldrsh struct {
	base

	rt, rn string
	off    int64
}

// newLdrsh - the Ldrsh constructor: the struct is assembled only
// here (the Builder method and the decoder call it).
func newLdrsh(b base, rt string, rn string, off int64) Ldrsh {
	return Ldrsh{
		base: b,
		rt:   rt,
		rn:   rn,
		off:  off,
	}
}

const ldrshEnc uint32 = 0x79800000

func (i Ldrsh) ObjDump(_ disasm.ViewCtx) string {
	if i.off == 0 {
		return fmt.Sprintf("ldrsh %s, [%s]", i.rt, i.rn)
	}

	return fmt.Sprintf("ldrsh %s, [%s, #0x%x]", i.rt, i.rn, i.off)
}

func (i Ldrsh) Encode(w io.Writer) (int64, error) {
	return lsSignedWrite(w, ldrshEnc, i.rt, i.rn, i.off, "ldrsh")
}

// Ldrsh — ldrsh rt, [rn, #off]: sign-extending halfword load, rt —
// x register only (register 31 reads as xzr), rn — x register or SP
// (register 31 in the base reads as sp); the offset is an imm12 scaled
// by 2 (0..0x1ffe, alignment 2).
func (Builder) Ldrsh(rt, rn Reg, off Off) (Instr, error) {
	if err := requireClass(rt, "Ldrsh", "rt", "x register (register 31 in rt reads as xzr)",
		classX, classXZR); err != nil {
		return nil, err
	}

	if err := requireClass(
		rn,
		"Ldrsh",
		"rn",
		"x register or SP (register 31 in the base reads as sp)",
		classX,
		classSP,
	); err != nil {
		return nil, err
	}

	if err := requireOff("Ldrsh", off, 1); err != nil {
		return nil, err
	}

	return newLdrsh(base{}, rt.name(), rn.name(), int64(off)), nil
}

func decodeLdrsh(w uint32) Instr {
	return newLdrsh(newBase(w), regNameX(w&0x1f), regNameXSP(w>>5&0x1f), int64(w>>10&0xfff)<<1)
}
