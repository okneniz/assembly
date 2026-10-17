package alias

// Generator for the csetm alias — one generator, one type, one text form
// family: csetm rd, cond.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"
	"github.com/okneniz/oh-snap/arbitrary"

	a64 "github.com/okneniz/assembly/arb/arm64"
	arm64 "github.com/okneniz/assembly/arch/arm64"
)

// CsetmParams — parameters of the csetm alias.
type CsetmParams struct {
	Rd   arm64.Reg
	Cond string
}

func NewCsetmParams(rd arm64.Reg, cond string) CsetmParams {
	return CsetmParams{
		Rd:   rd,
		Cond: cond,
	}
}

func (p CsetmParams) Instr() arm64.Instr {
	in, err := instrOfText(p.String())
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

func (p CsetmParams) String() string {
	return "csetm " + p.Rd.String() + ", " + p.Cond
}

// csetmGen — generator for csetm: a register and a condition (al/nv are
// not alias conditions).
type csetmGen struct {
	rnd  *rand.Rand
	coin ohsnap.Arbitrary[bool]
}

// Csetm — an arbitrary csetm.
func Csetm(rnd *rand.Rand) ohsnap.Arbitrary[CsetmParams] {
	return newCsetmGen(rnd)
}

func newCsetmGen(rnd *rand.Rand) csetmGen {
	return csetmGen{
		rnd:  rnd,
		coin: arbitrary.Bool(rnd),
	}
}

func (g csetmGen) Generate() CsetmParams {
	return NewCsetmParams(a64.GenReg(g.rnd, g.coin.Generate(), false, true), genCond14(g.rnd))
}

func (g csetmGen) Shrink(p CsetmParams) iter.Seq[CsetmParams] {
	var out []CsetmParams
	for _, r := range a64.RegShrunk(p.Rd) {
		out = append(out, NewCsetmParams(r, p.Cond))
	}

	if p.Cond != condCanonical() {
		out = append(out, NewCsetmParams(p.Rd, condCanonical()))
	}

	return slices.Values(out)
}
