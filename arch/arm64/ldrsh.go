package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Ldrsh — ldrsh rt, [rn, #imm12<<1].
type Ldrsh struct {
	rt, rn string
	off    int64
}

// newLdrsh - the Ldrsh constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newLdrsh(rt Reg, rn Reg, off Off) (Ldrsh, error) {
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
		rt:  rt.name(),
		rn:  rn.name(),
		off: int64(off),
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
