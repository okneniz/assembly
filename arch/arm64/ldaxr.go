package arm64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Ldaxr — ldaxr rt, [rn].
type Ldaxr struct {
	atomic

	enc uint32
}

// newLdaxrBase - the Ldaxr constructor for a ready embedded base
// (the decoder): the struct is assembled only here.
func newLdaxrBase(e atomic, enc uint32) Ldaxr {
	return Ldaxr{
		atomic: e,
		enc:    enc,
	}
}

// newLdaxr - the Ldaxr constructor: validates the operands,
// delegates the assembly to newLdaxrBase.
func newLdaxr(rt, rn Reg) (Ldaxr, error) {
	if err := lsOperand(rt, rn, "Ldaxr"); err != nil {
		return Ldaxr{}, err
	}

	enc := ldaxrWEnc
	if rt.Is64() {
		enc = ldaxrXEnc
	}

	return newLdaxrBase(newAtomic(rt.name(), rn.name()), enc), nil
}

// Encodings of the 64/32-bit forms: the access size is set by rt.
const (
	ldaxrXEnc uint32 = 0xC85FFC00 // ldaxr xt, [xn]
	ldaxrWEnc uint32 = 0x885FFC00 // ldaxr wt, [xn]
)

func (i Ldaxr) ObjDump(_ disasm.ViewCtx) string {
	return "ldaxr " + i.atText()
}

func (i Ldaxr) Encode(w io.Writer) (int64, error) {
	return i.atWrite(w, i.enc, "ldaxr")
}
