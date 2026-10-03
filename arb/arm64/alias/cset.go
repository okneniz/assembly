package alias

// Generator for the cset alias — one generator, one type, one text form
// family: cset rd, cond.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	a64 "github.com/okneniz/assembly/arb/arm64"
	arm64 "github.com/okneniz/assembly/arch/arm64"
)

// CsetParams — parameters of the cset alias.
type CsetParams struct {
	Rd   arm64.Reg
	Cond string
}

func NewCsetParams(rd arm64.Reg, cond string) CsetParams {
	return CsetParams{
		Rd:   rd,
		Cond: cond,
	}
}

func (p CsetParams) Instr() arm64.Instr {
	in, err := instrOfText(p.String())
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

func (p CsetParams) String() string {
	return "cset " + p.Rd.String() + ", " + p.Cond
}

// csetGen — generator for cset: a register and a condition (al/nv are
// not alias conditions).
type csetGen struct {
	rnd *rand.Rand
}

// Cset — an arbitrary cset.
func Cset(rnd *rand.Rand) ohsnap.Arbitrary[CsetParams] {
	return newCsetGen(rnd)
}

func newCsetGen(rnd *rand.Rand) csetGen {
	return csetGen{rnd: rnd}
}

func (g csetGen) Generate() iter.Seq[CsetParams] {
	return stream(func() CsetParams {
		return NewCsetParams(a64.GenReg(g.rnd, g.rnd.IntN(2) == 1, false, true), genCond14(g.rnd))
	})
}

func (g csetGen) Shrink(p CsetParams) iter.Seq[CsetParams] {
	var out []CsetParams
	for _, r := range a64.RegShrunk(p.Rd) {
		out = append(out, NewCsetParams(r, p.Cond))
	}

	if p.Cond != condCanonical() {
		out = append(out, NewCsetParams(p.Rd, condCanonical()))
	}

	return slices.Values(out)
}
