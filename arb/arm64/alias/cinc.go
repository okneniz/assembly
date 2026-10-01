package alias

// Generator for the cinc alias — one generator, one type, one text form
// family: cinc rd, rn, cond.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	a64 "github.com/okneniz/assembly/arb/arm64"
	arm64 "github.com/okneniz/assembly/arch/arm64"
)

// CincParams — parameters of the cinc alias.
type CincParams struct {
	Rd, Rn arm64.Reg
	Cond   string
}

func NewCincParams(rd arm64.Reg, rn arm64.Reg, cond string) CincParams {
	return CincParams{
		Rd:   rd,
		Rn:   rn,
		Cond: cond,
	}
}

func (p CincParams) String() string {
	return "cinc " + p.Rd.String() + ", " + p.Rn.String() + ", " + p.Cond
}

func (p CincParams) Instr() arm64.Instr {
	in, err := instrOfText(p.String())
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

// cincGen — generator for cinc: same-width registers and a condition.
type cincGen struct {
	rnd *rand.Rand
}

func newCincGen(rnd *rand.Rand) cincGen {
	return cincGen{rnd: rnd}
}

// Cinc — an arbitrary cinc.
func Cinc(rnd *rand.Rand) ohsnap.Arbitrary[CincParams] {
	return newCincGen(rnd)
}

func (g cincGen) Generate() iter.Seq[CincParams] {
	return stream(func() CincParams {
		is64 := g.rnd.IntN(2) == 1
		return NewCincParams(
			a64.GenReg(g.rnd, is64, false, true),
			a64.GenReg(g.rnd, is64, false, true),
			genCond14(g.rnd),
		)
	})
}

func (g cincGen) Shrink(p CincParams) iter.Seq[CincParams] {
	var out []CincParams
	for _, r := range a64.RegShrunk(p.Rd) {
		out = append(out, NewCincParams(r, p.Rn, p.Cond))
	}

	for _, r := range a64.RegShrunk(p.Rn) {
		out = append(out, NewCincParams(p.Rd, r, p.Cond))
	}

	if p.Cond != condCanonical() {
		out = append(out, NewCincParams(p.Rd, p.Rn, condCanonical()))
	}

	return slices.Values(out)
}
