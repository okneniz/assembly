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

// newLdrsh - the Ldrsh constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newLdrsh(b base, rt Reg, rn Reg, off Off) (Ldrsh, error) {
	err := requireClass(
		rt,
		"Ldrsh",
		"rt",
		"x register (register 31 in rt reads as xzr)",
		classX,
		classXZR,
	)

	if err != nil {
		return Ldrsh{}, err
	}

	err = requireClass(
		rn,
		"Ldrsh",
		"rn",
		"x register or SP (register 31 in the base reads as sp)",
		classX,
		classSP,
	)

	if err != nil {
		return Ldrsh{}, err
	}

	if err = requireOff("Ldrsh", off, 1); err != nil {
		return Ldrsh{}, err
	}

	return Ldrsh{
		base: b,
		rt:   rt.name(),
		rn:   rn.name(),
		off:  int64(off),
	}, nil
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

func (Builder) Ldrsh(rt, rn Reg, off Off) (Instr, error) {
	return newLdrsh(base{}, rt, rn, off)
}

func decodeLdrsh(w uint32) (Instr, error) {
	in, err := newLdrsh(
		newBase(w),
		xOf(w&0x1f),
		xspOf(w>>5&0x1f),
		Off(int64(w>>10&0xfff)<<1),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}
