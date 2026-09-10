package arm64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Stlr — stlr rt, [rn].
type Stlr struct {
	base
	atomic

	enc uint32
}

// newStlrBase - the Stlr constructor for a ready embedded base
// (the decoder): the struct is assembled only here.
func newStlrBase(b base, e atomic, enc uint32) Stlr {
	return Stlr{
		base:   b,
		atomic: e,
		enc:    enc,
	}
}

// newStlr - the Stlr constructor: validates the operands,
// delegates the assembly to newStlrBase.
func newStlr(b base, rt, rn Reg) (Stlr, error) {
	if err := lsOperand(rt, rn, "Stlr"); err != nil {
		return Stlr{}, err
	}

	enc := stlrWEnc
	if rt.Is64() {
		enc = stlrXEnc
	}

	return newStlrBase(b, newAtomic(rt.name(), rn.name()), enc), nil
}

// Encodings of the 64/32-bit forms: the access size is set by rt.
const (
	stlrXEnc uint32 = 0xC89FFC00 // stlr xt, [xn]
	stlrWEnc uint32 = 0x889FFC00 // stlr wt, [xn]
)

func decodeStlrOf(enc uint32, x64 bool) func(uint32) (Instr, error) {
	return func(w uint32) (Instr, error) {
		return Stlr{
			base:   newBase(w),
			atomic: newAtomic(armRegName(w&0x1f, x64), regNameXSP(w>>5&0x1f)),
			enc:    enc,
		}, nil
	}
}

func (i Stlr) ObjDump(_ disasm.ViewCtx) string {
	return "stlr " + i.atText()
}

func (i Stlr) Encode(w io.Writer) (int64, error) {
	return i.atWrite(w, i.enc, "stlr")
}

func (Builder) Stlr(rt, rn Reg) (Instr, error) {
	return newStlr(base{}, rt, rn)
}
