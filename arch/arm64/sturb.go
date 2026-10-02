package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Sturb — sturb ... (see lsBase for addressing kinds).
type Sturb struct {
	lsBase
}

// newSturbBase - the Sturb constructor for a ready embedded base
// (the decoder and the string-operand layer): the struct is
// assembled only here.
func newSturbBase(e lsBase) Sturb {
	return Sturb{
		lsBase: e,
	}
}

// newSturb - the Sturb constructor: validates the operands,
// delegates the assembly to newSturbBase.
func newSturb(rt, rn Reg, off Off) (Sturb, error) {
	err := requireClass(
		rt,
		"Sturb",
		"rt",
		"w register (register 31 in rt reads as wzr)",
		classW,
		classWZR,
	)

	if err != nil {
		return Sturb{}, err
	}

	err = requireClass(
		rn,
		"Sturb",
		"rn",
		"x register or SP (register 31 in the base reads as sp)",
		classX,
		classSP,
	)

	if err != nil {
		return Sturb{}, err
	}

	if err = requireUnscaledOff("Sturb", off); err != nil {
		return Sturb{}, err
	}

	return newSturbBase(
		newLsBase(rt.name(), rn.name(), memUnscaled, int64(off), 0, sturbEnc, "", "", 0),
	), nil
}

const sturbEnc uint32 = 0x38000000 // sturb wt, [xn, #±imm9]

func (i Sturb) ObjDump(ctx disasm.ViewCtx) string {
	return fmt.Sprintf("sturb %s, %s", i.rt, i.lsText(ctx))
}

func (i Sturb) Encode(w io.Writer) (int64, error) {
	return i.lsWrite(w, "sturb")
}
