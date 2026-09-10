package arm64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Stlrb — stlrb rt, [rn].
type Stlrb struct {
	base
	atomic

	enc uint32
}

// newStlrbBase - the Stlrb constructor for a ready embedded base
// (the decoder): the struct is assembled only here.
func newStlrbBase(b base, e atomic, enc uint32) Stlrb {
	return Stlrb{
		base:   b,
		atomic: e,
		enc:    enc,
	}
}

// newStlrb - the Stlrb constructor: validates the operands,
// delegates the assembly to newStlrbBase.
func newStlrb(b base, rt, rn Reg) (Stlrb, error) {
	err := requireClass(
		rt,
		"Stlrb",
		"rt",
		"w register (register 31 in rt reads as wzr)",
		classW,
		classWZR,
	)

	if err != nil {
		return Stlrb{}, err
	}

	err = requireClass(
		rn,
		"Stlrb",
		"rn",
		"x register or SP (register 31 in the base reads as sp)",
		classX,
		classSP,
	)

	if err != nil {
		return Stlrb{}, err
	}

	return newStlrbBase(b, newAtomic(rt.name(), rn.name()), stlrbEnc), nil
}

const stlrbEnc uint32 = 0x089FFC00 // stlrb wt, [xn]

func decodeStlrbOf(enc uint32, x64 bool) func(uint32) (Instr, error) {
	return func(w uint32) (Instr, error) {
		return Stlrb{
			base:   newBase(w),
			atomic: newAtomic(armRegName(w&0x1f, x64), regNameXSP(w>>5&0x1f)),
			enc:    enc,
		}, nil
	}
}

func (i Stlrb) ObjDump(_ disasm.ViewCtx) string {
	return "stlrb " + i.atText()
}

func (i Stlrb) Encode(w io.Writer) (int64, error) {
	return i.atWrite(w, i.enc, "stlrb")
}

func (Builder) Stlrb(rt, rn Reg) (Instr, error) {
	return newStlrb(base{}, rt, rn)
}
