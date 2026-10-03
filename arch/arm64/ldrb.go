package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Ldrb — ldrb ... (see lsBase for the addressing kinds).
type Ldrb struct {
	lsBase
}

// newLdrbBase - the Ldrb constructor for a ready embedded base
// (the decoder and the string-operand layer): the struct is
// assembled only here.
func newLdrbBase(e lsBase) Ldrb {
	return Ldrb{
		lsBase: e,
	}
}

// newLdrb - the Ldrb constructor: validates the operands,
// delegates the assembly to newLdrbBase.
func newLdrb(rt, rn Reg, off Off) (Ldrb, error) {
	err := requireClass(
		rt,
		"Ldrb",
		"rt",
		"w register (register 31 in rt reads as wzr)",
		classW,
		classWZR,
	)

	if err != nil {
		return Ldrb{}, err
	}

	err = requireClass(
		rn,
		"Ldrb",
		"rn",
		"x register or SP (register 31 in the base reads as sp)",
		classX,
		classSP,
	)

	if err != nil {
		return Ldrb{}, err
	}

	if err = requireOff("Ldrb", off, 0); err != nil {
		return Ldrb{}, err
	}

	return newLdrbBase(
		newLsBase(rt.name(), rn.name(), memImm, int64(off), 0, ldrbEnc, "", "", 0),
	), nil
}

const ldrbEnc uint32 = 0x39400000 // ldrb wt, [xn, #imm12]

func (i Ldrb) Encode(w io.Writer) (int64, error) {
	return i.lsWrite(w, "ldrb")
}

func (i Ldrb) ObjDump(ctx disasm.ViewCtx) string {
	return fmt.Sprintf("ldrb %s, %s", i.rt, i.lsText(ctx))
}
