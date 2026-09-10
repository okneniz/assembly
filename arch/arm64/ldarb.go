package arm64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Ldarb — ldarb rt, [rn].
type Ldarb struct {
	base
	atomic

	enc uint32
}

// newLdarbBase - the Ldarb constructor for a ready embedded base
// (the decoder): the struct is assembled only here.
func newLdarbBase(b base, e atomic, enc uint32) Ldarb {
	return Ldarb{
		base:   b,
		atomic: e,
		enc:    enc,
	}
}

// newLdarb - the Ldarb constructor: validates the operands,
// delegates the assembly to newLdarbBase.
func newLdarb(b base, rt, rn Reg) (Ldarb, error) {
	err := requireClass(
		rt,
		"Ldarb",
		"rt",
		"w register (register 31 in rt reads as wzr)",
		classW,
		classWZR,
	)

	if err != nil {
		return Ldarb{}, err
	}

	err = requireClass(
		rn,
		"Ldarb",
		"rn",
		"x register or SP (register 31 in the base reads as sp)",
		classX,
		classSP,
	)

	if err != nil {
		return Ldarb{}, err
	}

	return newLdarbBase(b, newAtomic(rt.name(), rn.name()), ldarbEnc), nil
}

const ldarbEnc uint32 = 0x08DFFC00 // ldarb wt, [xn]

func decodeLdarbOf(enc uint32, x64 bool) func(uint32) (Instr, error) {
	return func(w uint32) (Instr, error) {
		return Ldarb{
			base:   newBase(w),
			atomic: newAtomic(armRegName(w&0x1f, x64), regNameXSP(w>>5&0x1f)),
			enc:    enc,
		}, nil
	}
}

func (i Ldarb) ObjDump(_ disasm.ViewCtx) string {
	return "ldarb " + i.atText()
}

func (i Ldarb) Encode(w io.Writer) (int64, error) {
	return i.atWrite(w, i.enc, "ldarb")
}

func (Builder) Ldarb(rt, rn Reg) (Instr, error) {
	return newLdarb(base{}, rt, rn)
}
