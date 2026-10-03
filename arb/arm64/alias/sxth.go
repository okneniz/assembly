package alias

// Generator for the sxth alias — one generator, one type, one text form
// family: sxth xd|wd, wn (both destination widths).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	a64 "github.com/okneniz/assembly/arb/arm64"
	arm64 "github.com/okneniz/assembly/arch/arm64"
)

// SxthParams — parameters of the sxth alias.
type SxthParams struct {
	Rd arm64.Reg // x- or w-register
	Rn arm64.Reg // w-register
}

func NewSxthParams(rd arm64.Reg, rn arm64.Reg) SxthParams {
	return SxthParams{
		Rd: rd,
		Rn: rn,
	}
}

func (p SxthParams) Instr() arm64.Instr {
	in, err := instrOfText(p.String())
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

func (p SxthParams) String() string {
	return "sxth " + p.Rd.String() + ", " + p.Rn.String()
}

// sxthGen — generator for sxth: the destination of either width, the
// source a w-register.
type sxthGen struct {
	rnd *rand.Rand
}

// Sxth — an arbitrary sxth.
func Sxth(rnd *rand.Rand) ohsnap.Arbitrary[SxthParams] {
	return newSxthGen(rnd)
}

func newSxthGen(rnd *rand.Rand) sxthGen {
	return sxthGen{rnd: rnd}
}

func (g sxthGen) Generate() iter.Seq[SxthParams] {
	return stream(func() SxthParams {
		return NewSxthParams(
			a64.GenReg(g.rnd, g.rnd.IntN(2) == 1, false, true),
			a64.GenReg(g.rnd, false, false, true),
		)
	})
}

func (g sxthGen) Shrink(p SxthParams) iter.Seq[SxthParams] {
	regs := a64.RegShrunk(p.Rd)
	out := make([]SxthParams, 0, len(regs))
	for _, r := range regs {
		out = append(out, NewSxthParams(r, p.Rn))
	}

	for _, r := range a64.RegShrunk(p.Rn) {
		out = append(out, NewSxthParams(p.Rd, r))
	}

	return slices.Values(out)
}
