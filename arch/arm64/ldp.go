package arm64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Ldp — ldp rt, rt2, [rn{, #imm7<<scale}{!}].
type Ldp struct {
	base
	pairBase
}

// newLdpBase - the Ldp constructor for a ready embedded base
// (the decoder and the string-operand layer): the struct is
// assembled only here.
func newLdpBase(b base, e pairBase) Ldp {
	return Ldp{
		base:     b,
		pairBase: e,
	}
}

// newLdp - the Ldp constructor: validates the operands,
// delegates the assembly to newLdpBase.
func newLdp(b base, rt, rt2, rn Reg, off Off) (Ldp, error) {
	err := requireClass(
		rt,
		"Ldp",
		"rt",
		"x/w register (register 31 in rt reads as zr)",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return Ldp{}, err
	}

	err = requireClass(
		rt2,
		"Ldp",
		"rt2",
		"x/w register (register 31 in rt2 reads as zr)",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return Ldp{}, err
	}

	err = requireClass(
		rn,
		"Ldp",
		"rn",
		"x register or SP (register 31 in the base reads as sp)",
		classX,
		classSP,
	)

	if err != nil {
		return Ldp{}, err
	}

	err = requireWidth(
		"Ldp",
		rt,
		rt2,
	)

	if err != nil {
		return Ldp{}, err
	}

	enc, scale := ldpXEnc, uint32(3)
	if !rt.Is64() {
		enc, scale = ldpWEnc, 2
	}

	if err = requirePairOff("Ldp", off, scale); err != nil {
		return Ldp{}, err
	}

	return newLdpBase(
		b,
		newPairBase(rt.name(), rt2.name(), rn.name(), memImm, int64(off), scale, enc),
	), nil
}

// Encodings of the signed-offset form: the access size is set by rt.
const (
	ldpXEnc uint32 = 0xA9400000 // ldp xt, xt2, [xn, #imm7<<3]
	ldpWEnc uint32 = 0x29400000 // ldp wt, wt2, [xn, #imm7<<2]
)

func decodeLdpOf(enc uint32, scale uint32, x64 bool, rtKind string) func(uint32) (Instr, error) {
	kind := rtKind
	if kind == "" {
		if x64 {
			kind = "x"
		} else {
			kind = "w"
		}
	}

	return func(w uint32) (Instr, error) {
		rt, rt2, rn, k, off, load := pairDecode(w, scale, kind)
		if !load {
			return newStpBase(
				newBase(w),
				newPairBase(rt, rt2, rn, k, off, scale, enc&^1<<22),
			), nil
		}

		return Ldp{
			base:     newBase(w),
			pairBase: newPairBase(rt, rt2, rn, k, off, scale, enc|1<<22),
		}, nil
	}
}

func (i Ldp) ObjDump(_ disasm.ViewCtx) string {
	return "ldp " + i.pairText()
}

func (i Ldp) Encode(w io.Writer) (int64, error) {
	return i.pairWrite(w, "ldp")
}

func (Builder) Ldp(rt, rt2, rn Reg, off Off) (Instr, error) {
	return newLdp(base{}, rt, rt2, rn, off)
}
