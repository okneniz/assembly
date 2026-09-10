package arm64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Ldaxr — ldaxr rt, [rn].
type Ldaxr struct {
	base
	atomic

	enc uint32
}

// newLdaxrBase - the Ldaxr constructor for a ready embedded base
// (the decoder): the struct is assembled only here.
func newLdaxrBase(b base, e atomic, enc uint32) Ldaxr {
	return Ldaxr{
		base:   b,
		atomic: e,
		enc:    enc,
	}
}

// newLdaxr - the Ldaxr constructor: validates the operands,
// delegates the assembly to newLdaxrBase.
func newLdaxr(b base, rt, rn Reg) (Ldaxr, error) {
	if err := lsOperand(rt, rn, "Ldaxr"); err != nil {
		return Ldaxr{}, err
	}

	enc := ldaxrWEnc
	if rt.Is64() {
		enc = ldaxrXEnc
	}

	return newLdaxrBase(b, newAtomic(rt.name(), rn.name()), enc), nil
}

// Encodings of the 64/32-bit forms: the access size is set by rt.
const (
	ldaxrXEnc uint32 = 0xC85FFC00 // ldaxr xt, [xn]
	ldaxrWEnc uint32 = 0x885FFC00 // ldaxr wt, [xn]
)

func decodeLdaxrOf(enc uint32, x64 bool) func(uint32) (Instr, error) {
	return func(w uint32) (Instr, error) {
		return Ldaxr{
			base:   newBase(w),
			atomic: newAtomic(armRegName(w&0x1f, x64), regNameXSP(w>>5&0x1f)),
			enc:    enc,
		}, nil
	}
}

func (i Ldaxr) ObjDump(_ disasm.ViewCtx) string {
	return "ldaxr " + i.atText()
}

func (i Ldaxr) Encode(w io.Writer) (int64, error) {
	return i.atWrite(w, i.enc, "ldaxr")
}

// Ldaxr — ldaxr rt, [rn]: rt — x/w register (register 31 reads as
// zr), rn — x register or SP (register 31 in the base reads as sp).
func (Builder) Ldaxr(rt, rn Reg) (Instr, error) {
	return newLdaxr(base{}, rt, rn)
}
