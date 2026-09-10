package arm64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Stlxrb — stlxrb rs, rt, [rn].
type Stlxrb struct {
	base
	excl

	enc uint32
}

// newStlxrbBase - the Stlxrb constructor for a ready embedded base
// (the decoder): the struct is assembled only here.
func newStlxrbBase(b base, e excl, enc uint32) Stlxrb {
	return Stlxrb{
		base: b,
		excl: e,
		enc:  enc,
	}
}

// newStlxrb - the Stlxrb constructor: validates the operands,
// delegates the assembly to newStlxrbBase.
func newStlxrb(b base, rs, rt, rn Reg) (Stlxrb, error) {
	err := requireClass(
		rs,
		"Stlxrb",
		"rs",
		"w status register (register 31 in rs reads as wzr)",
		classW,
		classWZR,
	)

	if err != nil {
		return Stlxrb{}, err
	}

	err = requireClass(
		rt,
		"Stlxrb",
		"rt",
		"w register (register 31 in rt reads as wzr)",
		classW,
		classWZR,
	)

	if err != nil {
		return Stlxrb{}, err
	}

	err = requireClass(
		rn,
		"Stlxrb",
		"rn",
		"x register or SP (register 31 in the base reads as sp)",
		classX,
		classSP,
	)

	if err != nil {
		return Stlxrb{}, err
	}

	return newStlxrbBase(b, newExcl(rs.name(), rt.name(), rn.name()), stlxrbEnc), nil
}

const stlxrbEnc uint32 = 0x0800FC00 // stlxrb ws, wt, [xn]

func decodeStlxrbOf(enc uint32, x64 bool) func(uint32) (Instr, error) {
	return func(w uint32) (Instr, error) {
		return Stlxrb{
			base: newBase(w),
			excl: newExcl(regNameW(w>>16&0x1f), armRegName(w&0x1f, x64), regNameXSP(w>>5&0x1f)),
			enc:  enc,
		}, nil
	}
}

func (i Stlxrb) ObjDump(_ disasm.ViewCtx) string {
	return "stlxrb " + i.exText()
}

func (i Stlxrb) Encode(w io.Writer) (int64, error) {
	return i.exWrite(w, i.enc, "stlxrb")
}

func (Builder) Stlxrb(rs, rt, rn Reg) (Instr, error) {
	return newStlxrb(base{}, rs, rt, rn)
}
