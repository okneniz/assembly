package arm64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Stp — stp rt, rt2, [rn{, #imm7<<scale}{!}] (the same encoding, L=0).
type Stp struct {
	base
	pairBase
}

// newStpBase - the Stp constructor for a ready embedded base
// (the decoder and the string-operand layer): the struct is
// assembled only here.
func newStpBase(b base, e pairBase) Stp {
	return Stp{
		base:     b,
		pairBase: e,
	}
}

// newStp - the Stp constructor: validates the operands,
// delegates the assembly to newStpBase.
func newStp(b base, rt, rt2, rn Reg, off Off) (Stp, error) {
	err := requireClass(
		rt,
		"Stp",
		"rt",
		"x/w register (register 31 in rt reads as zr)",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return Stp{}, err
	}

	err = requireClass(
		rt2,
		"Stp",
		"rt2",
		"x/w register (register 31 in rt2 reads as zr)",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return Stp{}, err
	}

	err = requireClass(
		rn,
		"Stp",
		"rn",
		"x register or SP (register 31 in the base reads as sp)",
		classX,
		classSP,
	)

	if err != nil {
		return Stp{}, err
	}

	err = requireWidth(
		"Stp",
		rt,
		rt2,
	)

	if err != nil {
		return Stp{}, err
	}

	enc, scale := stpXEnc, uint32(3)
	if !rt.Is64() {
		enc, scale = stpWEnc, 2
	}

	if err = requirePairOff("Stp", off, scale); err != nil {
		return Stp{}, err
	}

	return newStpBase(
		b,
		newPairBase(rt.name(), rt2.name(), rn.name(), memImm, int64(off), scale, enc),
	), nil
}

// Encodings of the signed-offset form: the access size is set by rt.
const (
	stpXEnc uint32 = 0xA9000000 // stp xt, xt2, [xn, #imm7<<3]
	stpWEnc uint32 = 0x29000000 // stp wt, wt2, [xn, #imm7<<2]
)

func (i Stp) ObjDump(_ disasm.ViewCtx) string {
	return "stp " + i.pairText()
}

func (i Stp) Encode(w io.Writer) (int64, error) {
	return i.pairWrite(w, "stp")
}

func (Builder) Stp(rt, rt2, rn Reg, off Off) (Instr, error) {
	return newStp(base{}, rt, rt2, rn, off)
}
