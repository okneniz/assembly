package arm64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Ldar — ldar rt, [rn].
type Ldar struct {
	atomic

	enc uint32
}

// newLdarBase - the Ldar constructor for a ready embedded base
// (the decoder): the struct is assembled only here.
func newLdarBase(e atomic, enc uint32) Ldar {
	return Ldar{
		atomic: e,
		enc:    enc,
	}
}

// newLdar - the Ldar constructor: validates the operands,
// delegates the assembly to newLdarBase.
func newLdar(rt, rn Reg) (Ldar, error) {
	if err := lsOperand(rt, rn, "Ldar"); err != nil {
		return Ldar{}, err
	}

	enc := ldarWEnc
	if rt.Is64() {
		enc = ldarXEnc
	}

	return newLdarBase(newAtomic(rt.name(), rn.name()), enc), nil
}

// Encodings of the 64/32-bit forms: the access size is set by rt.
const (
	ldarXEnc uint32 = 0xC8DFFC00 // ldar xt, [xn]
	ldarWEnc uint32 = 0x88DFFC00 // ldar wt, [xn]
)

func (i Ldar) ObjDump(_ disasm.ViewCtx) string {
	return "ldar " + i.atText()
}

func (i Ldar) Encode(w io.Writer) (int64, error) {
	return i.atWrite(w, i.enc, "ldar")
}
