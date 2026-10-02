package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Ldrh — ldrh ... (see lsBase for the addressing kinds).
type Ldrh struct {
	base
	lsBase
}

// newLdrhBase - the Ldrh constructor for a ready embedded base
// (the decoder and the string-operand layer): the struct is
// assembled only here.
func newLdrhBase(b base, e lsBase) Ldrh {
	return Ldrh{
		base:   b,
		lsBase: e,
	}
}

// newLdrh - the Ldrh constructor: validates the operands,
// delegates the assembly to newLdrhBase.
func newLdrh(b base, rt, rn Reg, off Off) (Ldrh, error) {
	err := requireClass(
		rt,
		"Ldrh",
		"rt",
		"w register (register 31 in rt reads as wzr)",
		classW,
		classWZR,
	)

	if err != nil {
		return Ldrh{}, err
	}

	err = requireClass(
		rn,
		"Ldrh",
		"rn",
		"x register or SP (register 31 in the base reads as sp)",
		classX,
		classSP,
	)

	if err != nil {
		return Ldrh{}, err
	}

	if err = requireOff("Ldrh", off, 1); err != nil {
		return Ldrh{}, err
	}

	return newLdrhBase(
		b,
		newLsBase(rt.name(), rn.name(), memImm, int64(off), 0, ldrhEnc, "", "", 0),
	), nil
}

const ldrhEnc uint32 = 0x79400000 // ldrh wt, [xn, #imm12<<1]

func (i Ldrh) ObjDump(ctx disasm.ViewCtx) string {
	return fmt.Sprintf("ldrh %s, %s", i.rt, i.lsText(ctx))
}

func (i Ldrh) Encode(w io.Writer) (int64, error) {
	return i.lsWrite(w, "ldrh")
}
