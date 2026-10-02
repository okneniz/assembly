package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Ldurh — ldurh ... (see lsBase for the addressing kinds).
type Ldurh struct {
	lsBase
}

// newLdurhBase - the Ldurh constructor for a ready embedded base
// (the decoder and the string-operand layer): the struct is
// assembled only here.
func newLdurhBase(e lsBase) Ldurh {
	return Ldurh{
		lsBase: e,
	}
}

// newLdurh - the Ldurh constructor: validates the operands,
// delegates the assembly to newLdurhBase.
func newLdurh(rt, rn Reg, off Off) (Ldurh, error) {
	err := requireClass(
		rt,
		"Ldurh",
		"rt",
		"w register (register 31 in rt reads as wzr)",
		classW,
		classWZR,
	)

	if err != nil {
		return Ldurh{}, err
	}

	err = requireClass(
		rn,
		"Ldurh",
		"rn",
		"x register or SP (register 31 in the base reads as sp)",
		classX,
		classSP,
	)

	if err != nil {
		return Ldurh{}, err
	}

	if err = requireUnscaledOff("Ldurh", off); err != nil {
		return Ldurh{}, err
	}

	return newLdurhBase(
		newLsBase(rt.name(), rn.name(), memUnscaled, int64(off), 0, ldurhEnc, "", "", 0),
	), nil
}

const ldurhEnc uint32 = 0x78400000 // ldurh wt, [xn, #±imm9]

func (i Ldurh) ObjDump(ctx disasm.ViewCtx) string {
	return fmt.Sprintf("ldurh %s, %s", i.rt, i.lsText(ctx))
}

func (i Ldurh) Encode(w io.Writer) (int64, error) {
	return i.lsWrite(w, "ldurh")
}
