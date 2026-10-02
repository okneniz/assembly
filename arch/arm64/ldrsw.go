package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Ldrsw — ldrsw ... (see lsBase for the addressing kinds).
type Ldrsw struct {
	lsBase
}

// newLdrswBase - the Ldrsw constructor for a ready embedded base
// (the decoder and the string-operand layer): the struct is
// assembled only here.
func newLdrswBase(e lsBase) Ldrsw {
	return Ldrsw{
		lsBase: e,
	}
}

// newLdrsw - the Ldrsw constructor: validates the operands,
// delegates the assembly to newLdrswBase.
func newLdrsw(rt, rn Reg, off Off) (Ldrsw, error) {
	err := requireClass(
		rt,
		"Ldrsw",
		"rt",
		"x register (register 31 in rt reads as xzr)",
		classX,
		classXZR,
	)

	if err != nil {
		return Ldrsw{}, err
	}

	err = requireClass(
		rn,
		"Ldrsw",
		"rn",
		"x register or SP (register 31 in the base reads as sp)",
		classX,
		classSP,
	)

	if err != nil {
		return Ldrsw{}, err
	}

	if err = requireOff("Ldrsw", off, 2); err != nil {
		return Ldrsw{}, err
	}

	return newLdrswBase(
		newLsBase(rt.name(), rn.name(), memImm, int64(off), 0, ldrswEnc, "", "", 0),
	), nil
}

const ldrswEnc uint32 = 0xB9800000 // ldrsw xt, [xn, #imm12<<2]

func (i Ldrsw) ObjDump(ctx disasm.ViewCtx) string {
	return fmt.Sprintf("ldrsw %s, %s", i.rt, i.lsText(ctx))
}

func (i Ldrsw) Encode(w io.Writer) (int64, error) {
	return i.lsWrite(w, "ldrsw")
}
