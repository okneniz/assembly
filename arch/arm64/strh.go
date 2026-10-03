package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Strh — strh ... (see lsBase for addressing kinds).
type Strh struct {
	lsBase
}

// newStrhBase - the Strh constructor for a ready embedded base
// (the decoder and the string-operand layer): the struct is
// assembled only here.
func newStrhBase(e lsBase) Strh {
	return Strh{
		lsBase: e,
	}
}

// newStrh - the Strh constructor: validates the operands,
// delegates the assembly to newStrhBase.
func newStrh(rt, rn Reg, off Off) (Strh, error) {
	err := requireClass(
		rt,
		"Strh",
		"rt",
		"w register (register 31 in rt reads as wzr)",
		classW,
		classWZR,
	)

	if err != nil {
		return Strh{}, err
	}

	err = requireClass(
		rn,
		"Strh",
		"rn",
		"x register or SP (register 31 in the base reads as sp)",
		classX,
		classSP,
	)

	if err != nil {
		return Strh{}, err
	}

	if err = requireOff("Strh", off, 1); err != nil {
		return Strh{}, err
	}

	return newStrhBase(
		newLsBase(rt.name(), rn.name(), memImm, int64(off), 0, strhEnc, "", "", 0),
	), nil
}

const strhEnc uint32 = 0x79000000 // strh wt, [xn, #imm12<<1]

func (i Strh) Encode(w io.Writer) (int64, error) {
	return i.lsWrite(w, "strh")
}

func (i Strh) ObjDump(ctx disasm.ViewCtx) string {
	return fmt.Sprintf("strh %s, %s", i.rt, i.lsText(ctx))
}
