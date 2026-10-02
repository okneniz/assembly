package arm64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Stlr — stlr rt, [rn].
type Stlr struct {
	atomic

	enc uint32
}

// newStlrBase - the Stlr constructor for a ready embedded base
// (the decoder): the struct is assembled only here.
func newStlrBase(e atomic, enc uint32) Stlr {
	return Stlr{
		atomic: e,
		enc:    enc,
	}
}

// newStlr - the Stlr constructor: validates the operands,
// delegates the assembly to newStlrBase.
func newStlr(rt, rn Reg) (Stlr, error) {
	if err := lsOperand(rt, rn, "Stlr"); err != nil {
		return Stlr{}, err
	}

	enc := stlrWEnc
	if rt.Is64() {
		enc = stlrXEnc
	}

	return newStlrBase(newAtomic(rt.name(), rn.name()), enc), nil
}

// Encodings of the 64/32-bit forms: the access size is set by rt.
const (
	stlrXEnc uint32 = 0xC89FFC00 // stlr xt, [xn]
	stlrWEnc uint32 = 0x889FFC00 // stlr wt, [xn]
)

func (i Stlr) ObjDump(_ disasm.ViewCtx) string {
	return "stlr " + i.atText()
}

func (i Stlr) Encode(w io.Writer) (int64, error) {
	return i.atWrite(w, i.enc, "stlr")
}
