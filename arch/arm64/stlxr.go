package arm64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Stlxr — stlxr rs, rt, [rn].
type Stlxr struct {
	base
	excl

	enc uint32
}

// newStlxrBase - the Stlxr constructor for a ready embedded base
// (the decoder): the struct is assembled only here.
func newStlxrBase(b base, e excl, enc uint32) Stlxr {
	return Stlxr{
		base: b,
		excl: e,
		enc:  enc,
	}
}

// newStlxr - the Stlxr constructor: validates the operands,
// delegates the assembly to newStlxrBase.
func newStlxr(b base, rs, rt, rn Reg) (Stlxr, error) {
	err := requireClass(
		rs,
		"Stlxr",
		"rs",
		"w status register (register 31 in rs reads as wzr)",
		classW,
		classWZR,
	)

	if err != nil {
		return Stlxr{}, err
	}

	if err := lsOperand(rt, rn, "Stlxr"); err != nil {
		return Stlxr{}, err
	}

	enc := stlxrWEnc
	if rt.Is64() {
		enc = stlxrXEnc
	}

	return newStlxrBase(b, newExcl(rs.name(), rt.name(), rn.name()), enc), nil
}

// Encodings of the 64/32-bit forms: the access size is set by rt.
const (
	stlxrXEnc uint32 = 0xC800FC00 // stlxr ws, xt, [xn]
	stlxrWEnc uint32 = 0x8800FC00 // stlxr ws, wt, [xn]
)

func decodeStlxrOf(enc uint32, x64 bool) func(uint32) (Instr, error) {
	return func(w uint32) (Instr, error) {
		return Stlxr{
			base: newBase(w),
			excl: newExcl(regNameW(w>>16&0x1f), armRegName(w&0x1f, x64), regNameXSP(w>>5&0x1f)),
			enc:  enc,
		}, nil
	}
}

func (i Stlxr) ObjDump(_ disasm.ViewCtx) string {
	return "stlxr " + i.exText()
}

func (i Stlxr) Encode(w io.Writer) (int64, error) {
	return i.exWrite(w, i.enc, "stlxr")
}

func (Builder) Stlxr(rs, rt, rn Reg) (Instr, error) {
	return newStlxr(base{}, rs, rt, rn)
}
