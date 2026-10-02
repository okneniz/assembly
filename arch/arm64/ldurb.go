package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Ldurb — ldurb ... (see lsBase for the addressing kinds).
type Ldurb struct {
	base
	lsBase
}

// newLdurbBase - the Ldurb constructor for a ready embedded base
// (the decoder and the string-operand layer): the struct is
// assembled only here.
func newLdurbBase(b base, e lsBase) Ldurb {
	return Ldurb{
		base:   b,
		lsBase: e,
	}
}

// newLdurb - the Ldurb constructor: validates the operands,
// delegates the assembly to newLdurbBase.
func newLdurb(b base, rt, rn Reg, off Off) (Ldurb, error) {
	err := requireClass(
		rt,
		"Ldurb",
		"rt",
		"w register (register 31 in rt reads as wzr)",
		classW,
		classWZR,
	)

	if err != nil {
		return Ldurb{}, err
	}

	err = requireClass(
		rn,
		"Ldurb",
		"rn",
		"x register or SP (register 31 in the base reads as sp)",
		classX,
		classSP,
	)

	if err != nil {
		return Ldurb{}, err
	}

	if err = requireUnscaledOff("Ldurb", off); err != nil {
		return Ldurb{}, err
	}

	return newLdurbBase(
		b,
		newLsBase(rt.name(), rn.name(), memUnscaled, int64(off), 0, ldurbEnc, "", "", 0),
	), nil
}

const ldurbEnc uint32 = 0x38400000 // ldurb wt, [xn, #±imm9]

func (i Ldurb) ObjDump(ctx disasm.ViewCtx) string {
	return fmt.Sprintf("ldurb %s, %s", i.rt, i.lsText(ctx))
}

func (i Ldurb) Encode(w io.Writer) (int64, error) {
	return i.lsWrite(w, "ldurb")
}
