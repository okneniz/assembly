package alias

// Generator for the cneg alias — one generator, one type, one text form
// family: cneg rd, rn, cond.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	a64 "github.com/okneniz/assembly/arb/arm64"
	arm64 "github.com/okneniz/assembly/arch/arm64"
)

// CnegParams — parameters of the cneg alias.
type CnegParams struct {
	Rd, Rn arm64.Reg
	Cond   string
}

func NewCnegParams(rd arm64.Reg, rn arm64.Reg, cond string) CnegParams {
	return CnegParams{
		Rd:   rd,
		Rn:   rn,
		Cond: cond,
	}
}

func (p CnegParams) Instr() arm64.Instr {
	in, err := instrOfText(p.String())
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

func (p CnegParams) String() string {
	return "cneg " + p.Rd.String() + ", " + p.Rn.String() + ", " + p.Cond
}

// cnegGen — generator for cneg: same-width registers and a condition.
type cnegGen struct {
	rnd *rand.Rand
}

// Cneg — an arbitrary cneg.
func Cneg(rnd *rand.Rand) ohsnap.Arbitrary[CnegParams] {
	return newCnegGen(rnd)
}

func newCnegGen(rnd *rand.Rand) cnegGen {
	return cnegGen{rnd: rnd}
}

func (g cnegGen) Generate() iter.Seq[CnegParams] {
	return stream(func() CnegParams {
		is64 := g.rnd.IntN(2) == 1
		return NewCnegParams(
			a64.GenReg(g.rnd, is64, false, true),
			a64.GenReg(g.rnd, is64, false, true),
			genCond14(g.rnd),
		)
	})
}

func (g cnegGen) Shrink(p CnegParams) iter.Seq[CnegParams] {
	var out []CnegParams
	for _, r := range a64.RegShrunk(p.Rd) {
		out = append(out, NewCnegParams(r, p.Rn, p.Cond))
	}

	for _, r := range a64.RegShrunk(p.Rn) {
		out = append(out, NewCnegParams(p.Rd, r, p.Cond))
	}

	if p.Cond != condCanonical() {
		out = append(out, NewCnegParams(p.Rd, p.Rn, condCanonical()))
	}

	return slices.Values(out)
}
