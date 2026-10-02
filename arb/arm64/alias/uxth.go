package alias

// Generator for the uxth alias — one generator, one type, one text
// form family: uxth wd, wn (the UBFM #0, #15 encoding, a 32-bit form
// only).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	a64 "github.com/okneniz/assembly/arb/arm64"
	arm64 "github.com/okneniz/assembly/arch/arm64"
)

// UxthParams — parameters of the uxth alias.
type UxthParams struct {
	Rd, Rn arm64.Reg // w-registers
}

func NewUxthParams(rd arm64.Reg, rn arm64.Reg) UxthParams {
	return UxthParams{
		Rd: rd,
		Rn: rn,
	}
}

func (p UxthParams) String() string {
	return "uxth " + p.Rd.String() + ", " + p.Rn.String()
}

func (p UxthParams) Instr() arm64.Instr {
	in, err := instrOfText(p.String())
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

// uxthGen — generator for uxth: both registers are w-registers.
type uxthGen struct {
	rnd *rand.Rand
}

func newUxthGen(rnd *rand.Rand) uxthGen {
	return uxthGen{rnd: rnd}
}

// Uxth — an arbitrary uxth.
func Uxth(rnd *rand.Rand) ohsnap.Arbitrary[UxthParams] {
	return newUxthGen(rnd)
}

func (g uxthGen) Generate() iter.Seq[UxthParams] {
	return stream(func() UxthParams {
		return NewUxthParams(
			a64.GenReg(g.rnd, false, false, true),
			a64.GenReg(g.rnd, false, false, true),
		)
	})
}

func (g uxthGen) Shrink(p UxthParams) iter.Seq[UxthParams] {
	var out []UxthParams
	for _, r := range a64.RegShrunk(p.Rd) {
		out = append(out, NewUxthParams(r, p.Rn))
	}

	for _, r := range a64.RegShrunk(p.Rn) {
		out = append(out, NewUxthParams(p.Rd, r))
	}

	return slices.Values(out)
}
