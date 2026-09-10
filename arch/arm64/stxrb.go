package arm64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Stxrb — stxrb rs, rt, [rn].
type Stxrb struct {
	base
	excl

	enc uint32
}

// newStxrbBase - the Stxrb constructor for a ready embedded base
// (the decoder): the struct is assembled only here.
func newStxrbBase(b base, e excl, enc uint32) Stxrb {
	return Stxrb{
		base: b,
		excl: e,
		enc:  enc,
	}
}

// newStxrb - the Stxrb constructor: validates the operands,
// delegates the assembly to newStxrbBase.
func newStxrb(b base, rs, rt, rn Reg) (Stxrb, error) {
	err := requireClass(
		rs,
		"Stxrb",
		"rs",
		"w status register (register 31 in rs reads as wzr)",
		classW,
		classWZR,
	)

	if err != nil {
		return Stxrb{}, err
	}

	err = requireClass(
		rt,
		"Stxrb",
		"rt",
		"w register (register 31 in rt reads as wzr)",
		classW,
		classWZR,
	)

	if err != nil {
		return Stxrb{}, err
	}

	err = requireClass(
		rn,
		"Stxrb",
		"rn",
		"x register or SP (register 31 in the base reads as sp)",
		classX,
		classSP,
	)

	if err != nil {
		return Stxrb{}, err
	}

	return newStxrbBase(b, newExcl(rs.name(), rt.name(), rn.name()), stxrbEnc), nil
}

const stxrbEnc uint32 = 0x08000000 // stxrb ws, wt, [xn]

func decodeStxrbOf(enc uint32, x64 bool) func(uint32) (Instr, error) {
	return func(w uint32) (Instr, error) {
		return Stxrb{
			base: newBase(w),
			excl: newExcl(regNameW(w>>16&0x1f), armRegName(w&0x1f, x64), regNameXSP(w>>5&0x1f)),
			enc:  enc,
		}, nil
	}
}

func (i Stxrb) ObjDump(_ disasm.ViewCtx) string {
	return "stxrb " + i.exText()
}

func (i Stxrb) Encode(w io.Writer) (int64, error) {
	return i.exWrite(w, i.enc, "stxrb")
}

func (Builder) Stxrb(rs, rt, rn Reg) (Instr, error) {
	return newStxrb(base{}, rs, rt, rn)
}
