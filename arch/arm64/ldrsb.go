package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Ldrsb — ldrsb rt, [rn, #imm12] (Rt — X; the W form is not decoded).
type Ldrsb struct {
	rt, rn string
	off    int64
}

// newLdrsb - the Ldrsb constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newLdrsb(rt Reg, rn Reg, off Off) (Ldrsb, error) {
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
		rt:  rt.name(),
		rn:  rn.name(),
		off: int64(off),
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
