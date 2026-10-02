package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Strb — strb ... (see lsBase for addressing kinds).
type Strb struct {
	base
	lsBase
}

// newStrbBase - the Strb constructor for a ready embedded base (the
// decoder and the string-operand layer): the struct is assembled only
// here.
func newStrbBase(b base, e lsBase) Strb {
	return Strb{
		base:   b,
		lsBase: e,
	}
}

// newStrb - the Strb constructor: validates the operands,
// delegates the assembly to newStrbBase.
func newStrb(b base, rt, rn Reg, off Off) (Strb, error) {
	err := requireClass(
		rt,
		"Strb",
		"rt",
		"w register (register 31 in rt reads as wzr)",
		classW,
		classWZR,
	)

	if err != nil {
		return Strb{}, err
	}

	err = requireClass(
		rn,
		"Strb",
		"rn",
		"x register or SP (register 31 in the base reads as sp)",
		classX,
		classSP,
	)

	if err != nil {
		return Strb{}, err
	}

	if err = requireOff("Strb", off, 0); err != nil {
		return Strb{}, err
	}

	return newStrbBase(
		b,
		newLsBase(rt.name(), rn.name(), memImm, int64(off), 0, strbEnc, "", "", 0),
	), nil
}

const strbEnc uint32 = 0x39000000 // strb wt, [xn, #imm12]

func (i Strb) ObjDump(ctx disasm.ViewCtx) string {
	return fmt.Sprintf("strb %s, %s", i.rt, i.lsText(ctx))
}

func (i Strb) Encode(w io.Writer) (int64, error) {
	return i.lsWrite(w, "strb")
}
