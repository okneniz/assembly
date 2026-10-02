package arm64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Stlxrb — stlxrb rs, rt, [rn].
type Stlxrb struct {
	excl

	enc uint32
}

// newStlxrbBase - the Stlxrb constructor for a ready embedded base
// (the decoder): the struct is assembled only here.
func newStlxrbBase(e excl, enc uint32) Stlxrb {
	return Stlxrb{
		excl: e,
		enc:  enc,
	}
}

// newStlxrb - the Stlxrb constructor: validates the operands,
// delegates the assembly to newStlxrbBase.
func newStlxrb(rs, rt, rn Reg) (Stlxrb, error) {
	err := requireClass(
		rs,
		"Stlxrb",
		"rs",
		"w status register (register 31 in rs reads as wzr)",
		classW,
		classWZR,
	)

	if err != nil {
		return Stlxrb{}, err
	}

	err = requireClass(
		rt,
		"Stlxrb",
		"rt",
		"w register (register 31 in rt reads as wzr)",
		classW,
		classWZR,
	)

	if err != nil {
		return Stlxrb{}, err
	}

	err = requireClass(
		rn,
		"Stlxrb",
		"rn",
		"x register or SP (register 31 in the base reads as sp)",
		classX,
		classSP,
	)

	if err != nil {
		return Stlxrb{}, err
	}

	return newStlxrbBase(newExcl(rs.name(), rt.name(), rn.name()), stlxrbEnc), nil
}

const stlxrbEnc uint32 = 0x0800FC00 // stlxrb ws, wt, [xn]

func (i Stlxrb) ObjDump(_ disasm.ViewCtx) string {
	return "stlxrb " + i.exText()
}

func (i Stlxrb) Encode(w io.Writer) (int64, error) {
	return i.exWrite(w, i.enc, "stlxrb")
}
