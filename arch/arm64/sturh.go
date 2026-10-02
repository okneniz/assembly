package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Sturh — sturh ... (see lsBase for addressing kinds).
type Sturh struct {
	base
	lsBase
}

// newSturhBase - the Sturh constructor for a ready embedded base
// (the decoder and the string-operand layer): the struct is
// assembled only here.
func newSturhBase(b base, e lsBase) Sturh {
	return Sturh{
		base:   b,
		lsBase: e,
	}
}

// newSturh - the Sturh constructor: validates the operands,
// delegates the assembly to newSturhBase.
func newSturh(b base, rt, rn Reg, off Off) (Sturh, error) {
	err := requireClass(
		rt,
		"Sturh",
		"rt",
		"w register (register 31 in rt reads as wzr)",
		classW,
		classWZR,
	)

	if err != nil {
		return Sturh{}, err
	}

	err = requireClass(
		rn,
		"Sturh",
		"rn",
		"x register or SP (register 31 in the base reads as sp)",
		classX,
		classSP,
	)

	if err != nil {
		return Sturh{}, err
	}

	if err = requireUnscaledOff("Sturh", off); err != nil {
		return Sturh{}, err
	}

	return newSturhBase(
		b,
		newLsBase(rt.name(), rn.name(), memUnscaled, int64(off), 0, sturhEnc, "", "", 0),
	), nil
}

const sturhEnc uint32 = 0x78000000 // sturh wt, [xn, #±imm9]

func (i Sturh) ObjDump(ctx disasm.ViewCtx) string {
	return fmt.Sprintf("sturh %s, %s", i.rt, i.lsText(ctx))
}

func (i Sturh) Encode(w io.Writer) (int64, error) {
	return i.lsWrite(w, "sturh")
}
