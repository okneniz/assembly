package alias

// Generator for the sxtb alias — one generator, one type, one text form
// family: sxtb xd|wd, wn (both destination widths).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	a64 "github.com/okneniz/assembly/arb/arm64"
	arm64 "github.com/okneniz/assembly/arch/arm64"
)

// SxtbParams — parameters of the sxtb alias.
type SxtbParams struct {
	Rd arm64.Reg // x- or w-register
	Rn arm64.Reg // w-register
}

func NewSxtbParams(rd arm64.Reg, rn arm64.Reg) SxtbParams {
	return SxtbParams{
		Rd: rd,
		Rn: rn,
	}
}

func (p SxtbParams) Instr() arm64.Instr {
	in, err := instrOfText(p.String())
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

func (p SxtbParams) String() string {
	return "sxtb " + p.Rd.String() + ", " + p.Rn.String()
}

// sxtbGen — generator for sxtb: the destination of either width, the
// source a w-register.
type sxtbGen struct {
	rnd *rand.Rand
}

// Sxtb — an arbitrary sxtb.
func Sxtb(rnd *rand.Rand) ohsnap.Arbitrary[SxtbParams] {
	return newSxtbGen(rnd)
}

func newSxtbGen(rnd *rand.Rand) sxtbGen {
	return sxtbGen{rnd: rnd}
}

func (g sxtbGen) Generate() iter.Seq[SxtbParams] {
	return stream(func() SxtbParams {
		return NewSxtbParams(
			a64.GenReg(g.rnd, g.rnd.IntN(2) == 1, false, true),
			a64.GenReg(g.rnd, false, false, true),
		)
	})
}

func (g sxtbGen) Shrink(p SxtbParams) iter.Seq[SxtbParams] {
	regs := a64.RegShrunk(p.Rd)
	out := make([]SxtbParams, 0, len(regs))
	for _, r := range regs {
		out = append(out, NewSxtbParams(r, p.Rn))
	}

	for _, r := range a64.RegShrunk(p.Rn) {
		out = append(out, NewSxtbParams(p.Rd, r))
	}

	return slices.Values(out)
}
