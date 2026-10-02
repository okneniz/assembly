package arm64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Ldaxrb — ldaxrb rt, [rn].
type Ldaxrb struct {
	atomic

	enc uint32
}

// newLdaxrbBase - the Ldaxrb constructor for a ready embedded base
// (the decoder): the struct is assembled only here.
func newLdaxrbBase(e atomic, enc uint32) Ldaxrb {
	return Ldaxrb{
		atomic: e,
		enc:    enc,
	}
}

// newLdaxrb - the Ldaxrb constructor: validates the operands,
// delegates the assembly to newLdaxrbBase.
func newLdaxrb(rt, rn Reg) (Ldaxrb, error) {
	err := requireClass(
		rt,
		"Ldaxrb",
		"rt",
		"w register (register 31 in rt reads as wzr)",
		classW,
		classWZR,
	)

	if err != nil {
		return Ldaxrb{}, err
	}

	err = requireClass(
		rn,
		"Ldaxrb",
		"rn",
		"x register or SP (register 31 in the base reads as sp)",
		classX,
		classSP,
	)

	if err != nil {
		return Ldaxrb{}, err
	}

	return newLdaxrbBase(newAtomic(rt.name(), rn.name()), ldaxrbEnc), nil
}

const ldaxrbEnc uint32 = 0x085FFC00 // ldaxrb wt, [xn]

func (i Ldaxrb) ObjDump(_ disasm.ViewCtx) string {
	return "ldaxrb " + i.atText()
}

func (i Ldaxrb) Encode(w io.Writer) (int64, error) {
	return i.atWrite(w, i.enc, "ldaxrb")
}
