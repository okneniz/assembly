package alias

// Generator for the uxtb alias — one generator, one type, one text
// form family: uxtb wd, wn (the UBFM #0, #7 encoding, a 32-bit form
// only).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	a64 "github.com/okneniz/assembly/arb/arm64"
	arm64 "github.com/okneniz/assembly/arch/arm64"
)

// UxtbParams — parameters of the uxtb alias.
type UxtbParams struct {
	Rd, Rn arm64.Reg // w-registers
}

func NewUxtbParams(rd arm64.Reg, rn arm64.Reg) UxtbParams {
	return UxtbParams{
		Rd: rd,
		Rn: rn,
	}
}

func (p UxtbParams) Instr() arm64.Instr {
	in, err := instrOfText(p.String())
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

func (p UxtbParams) String() string {
	return "uxtb " + p.Rd.String() + ", " + p.Rn.String()
}

// uxtbGen — generator for uxtb: both registers are w-registers.
type uxtbGen struct {
	rnd *rand.Rand
}

// Uxtb — an arbitrary uxtb.
func Uxtb(rnd *rand.Rand) ohsnap.Arbitrary[UxtbParams] {
	return newUxtbGen(rnd)
}

func newUxtbGen(rnd *rand.Rand) uxtbGen {
	return uxtbGen{rnd: rnd}
}

func (g uxtbGen) Generate() iter.Seq[UxtbParams] {
	return stream(func() UxtbParams {
		return NewUxtbParams(
			a64.GenReg(g.rnd, false, false, true),
			a64.GenReg(g.rnd, false, false, true),
		)
	})
}

func (g uxtbGen) Shrink(p UxtbParams) iter.Seq[UxtbParams] {
	regs := a64.RegShrunk(p.Rd)
	out := make([]UxtbParams, 0, len(regs))
	for _, r := range regs {
		out = append(out, NewUxtbParams(r, p.Rn))
	}

	for _, r := range a64.RegShrunk(p.Rn) {
		out = append(out, NewUxtbParams(p.Rd, r))
	}

	return slices.Values(out)
}
