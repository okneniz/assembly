package alias

// Generator for the sxtw alias — one generator, one type, one text form
// family: sxtw xd, wn.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	a64 "github.com/okneniz/assembly/arb/arm64"
	arm64 "github.com/okneniz/assembly/arch/arm64"
)

// SxtwParams — parameters of the sxtw alias.
type SxtwParams struct {
	Rd arm64.Reg // x-register
	Rn arm64.Reg // w-register
}

func NewSxtwParams(rd arm64.Reg, rn arm64.Reg) SxtwParams {
	return SxtwParams{
		Rd: rd,
		Rn: rn,
	}
}

func (p SxtwParams) String() string {
	return "sxtw " + p.Rd.String() + ", " + p.Rn.String()
}

func (p SxtwParams) Instr() arm64.Instr {
	in, err := instrOfText(p.String())
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

// sxtwGen — generator for sxtw: the destination is an x-register, the
// source a w-register.
type sxtwGen struct {
	rnd *rand.Rand
}

func newSxtwGen(rnd *rand.Rand) sxtwGen {
	return sxtwGen{rnd: rnd}
}

// Sxtw — an arbitrary sxtw.
func Sxtw(rnd *rand.Rand) ohsnap.Arbitrary[SxtwParams] {
	return newSxtwGen(rnd)
}

func (g sxtwGen) Generate() iter.Seq[SxtwParams] {
	return stream(func() SxtwParams {
		return NewSxtwParams(a64.GenReg(g.rnd, true, false, true), a64.GenReg(g.rnd, false, false, true))
	})
}

func (g sxtwGen) Shrink(p SxtwParams) iter.Seq[SxtwParams] {
	var out []SxtwParams
	for _, r := range a64.RegShrunk(p.Rd) {
		out = append(out, NewSxtwParams(r, p.Rn))
	}

	for _, r := range a64.RegShrunk(p.Rn) {
		out = append(out, NewSxtwParams(p.Rd, r))
	}

	return slices.Values(out)
}
