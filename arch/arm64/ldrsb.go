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

// newLdrsb - the Ldrsb constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newLdrsb(b base, rt Reg, rn Reg, off Off) (Ldrsb, error) {
	err := requireClass(
		rt,
		"Ldrsb",
		"rt",
		"x register (register 31 in rt reads as xzr)",
		classX,
		classXZR,
	)

	if err != nil {
		return Ldrsb{}, err
	}

	err = requireClass(
		rn,
		"Ldrsb",
		"rn",
		"x register or SP (register 31 in the base reads as sp)",
		classX,
		classSP,
	)

	if err != nil {
		return Ldrsb{}, err
	}

	if err = requireOff("Ldrsb", off, 0); err != nil {
		return Ldrsb{}, err
	}

	return Ldrsb{
		base: b,
		rt:   rt.name(),
		rn:   rn.name(),
		off:  int64(off),
	}, nil
}

const ldrsbEnc uint32 = 0x39800000

func (i Ldrsb) ObjDump(_ disasm.ViewCtx) string {
	if i.off == 0 {
		return fmt.Sprintf("ldrsb %s, [%s]", i.rt, i.rn)
	}

	return fmt.Sprintf("ldrsb %s, [%s, #0x%x]", i.rt, i.rn, i.off)
}

func (i Ldrsb) Encode(w io.Writer) (int64, error) {
	return lsSignedWrite(w, ldrsbEnc, i.rt, i.rn, i.off, "ldrsb")
}

func (Builder) Ldrsb(rt, rn Reg, off Off) (Instr, error) {
	return newLdrsb(base{}, rt, rn, off)
}

func decodeLdrsb(w uint32) (Instr, error) {
	in, err := newLdrsb(
		newBase(w),
		xOf(w&0x1f),
		xspOf(w>>5&0x1f),
		Off(int64(w>>10&0xfff)<<0),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}
