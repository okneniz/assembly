package alias

// Generator for the cinv alias — one generator, one type, one text form
// family: cinv rd, rn, cond.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	a64 "github.com/okneniz/assembly/arb/arm64"
	arm64 "github.com/okneniz/assembly/arch/arm64"
)

// CinvParams — parameters of the cinv alias.
type CinvParams struct {
	Rd, Rn arm64.Reg
	Cond   string
}

func NewCinvParams(rd arm64.Reg, rn arm64.Reg, cond string) CinvParams {
	return CinvParams{
		Rd:   rd,
		Rn:   rn,
		Cond: cond,
	}
}

func (p CinvParams) Instr() arm64.Instr {
	in, err := instrOfText(p.String())
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

func (p CinvParams) String() string {
	return "cinv " + p.Rd.String() + ", " + p.Rn.String() + ", " + p.Cond
}

// cinvGen — generator for cinv: same-width registers and a condition.
type cinvGen struct {
	rnd *rand.Rand
}

// Cinv — an arbitrary cinv.
func Cinv(rnd *rand.Rand) ohsnap.Arbitrary[CinvParams] {
	return newCinvGen(rnd)
}

func newCinvGen(rnd *rand.Rand) cinvGen {
	return cinvGen{rnd: rnd}
}

func (g cinvGen) Generate() iter.Seq[CinvParams] {
	return stream(func() CinvParams {
		is64 := g.rnd.IntN(2) == 1
		return NewCinvParams(
			a64.GenReg(g.rnd, is64, false, true),
			a64.GenReg(g.rnd, is64, false, true),
			genCond14(g.rnd),
		)
	})
}

func (g cinvGen) Shrink(p CinvParams) iter.Seq[CinvParams] {
	var out []CinvParams
	for _, r := range a64.RegShrunk(p.Rd) {
		out = append(out, NewCinvParams(r, p.Rn, p.Cond))
	}

	for _, r := range a64.RegShrunk(p.Rn) {
		out = append(out, NewCinvParams(p.Rd, r, p.Cond))
	}

	if p.Cond != condCanonical() {
		out = append(out, NewCinvParams(p.Rd, p.Rn, condCanonical()))
	}

	return slices.Values(out)
}
