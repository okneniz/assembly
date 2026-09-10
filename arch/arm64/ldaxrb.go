package arm64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Ldaxrb — ldaxrb rt, [rn].
type Ldaxrb struct {
	base
	atomic

	enc uint32
}

// newLdaxrbBase - the Ldaxrb constructor for a ready embedded base
// (the decoder): the struct is assembled only here.
func newLdaxrbBase(b base, e atomic, enc uint32) Ldaxrb {
	return Ldaxrb{
		base:   b,
		atomic: e,
		enc:    enc,
	}
}

// newLdaxrb - the Ldaxrb constructor: validates the operands,
// delegates the assembly to newLdaxrbBase.
func newLdaxrb(b base, rt, rn Reg) (Ldaxrb, error) {
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

	return newLdaxrbBase(b, newAtomic(rt.name(), rn.name()), ldaxrbEnc), nil
}

const ldaxrbEnc uint32 = 0x085FFC00 // ldaxrb wt, [xn]

func decodeLdaxrbOf(enc uint32, x64 bool) func(uint32) (Instr, error) {
	return func(w uint32) (Instr, error) {
		return Ldaxrb{
			base:   newBase(w),
			atomic: newAtomic(armRegName(w&0x1f, x64), regNameXSP(w>>5&0x1f)),
			enc:    enc,
		}, nil
	}
}

func (i Ldaxrb) ObjDump(_ disasm.ViewCtx) string {
	return "ldaxrb " + i.atText()
}

func (i Ldaxrb) Encode(w io.Writer) (int64, error) {
	return i.atWrite(w, i.enc, "ldaxrb")
}

// Ldaxrb — ldaxrb rt, [rn]: byte access, rt — w register only
// (register 31 reads as wzr), rn — x register or SP (register 31 in the
// base reads as sp).
func (Builder) Ldaxrb(rt, rn Reg) (Instr, error) {
	return newLdaxrb(base{}, rt, rn)
}
