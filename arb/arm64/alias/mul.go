package alias

// Generator for the mul alias — one generator, one type, one text form
// family: mul rd, rn, rm.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	a64 "github.com/okneniz/assembly/arb/arm64"
	arm64 "github.com/okneniz/assembly/arch/arm64"
)

// MulParams — parameters of the mul alias.
type MulParams struct {
	Rd, Rn, Rm arm64.Reg
}

func NewMulParams(rd arm64.Reg, rn arm64.Reg, rm arm64.Reg) MulParams {
	return MulParams{
		Rd: rd,
		Rn: rn,
		Rm: rm,
	}
}

func (p MulParams) Instr() arm64.Instr {
	in, err := instrOfText(p.String())
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

func (p MulParams) String() string {
	return "mul " + p.Rd.String() + ", " + p.Rn.String() + ", " + p.Rm.String()
}

// mulGen — generator for mul: same-width registers, 31st is zr.
type mulGen struct {
	rnd *rand.Rand
}

// Mul — an arbitrary mul.
func Mul(rnd *rand.Rand) ohsnap.Arbitrary[MulParams] {
	return newMulGen(rnd)
}

func newMulGen(rnd *rand.Rand) mulGen {
	return mulGen{rnd: rnd}
}

func (g mulGen) Generate() iter.Seq[MulParams] {
	return stream(func() MulParams {
		is64 := g.rnd.IntN(2) == 1
		return NewMulParams(
			a64.GenReg(g.rnd, is64, false, true),
			a64.GenReg(g.rnd, is64, false, true),
			a64.GenReg(g.rnd, is64, false, true),
		)
	})
}

func (g mulGen) Shrink(p MulParams) iter.Seq[MulParams] {
	regs := a64.RegShrunk(p.Rd)
	out := make([]MulParams, 0, len(regs))
	for _, r := range regs {
		out = append(out, NewMulParams(r, p.Rn, p.Rm))
	}

	for _, r := range a64.RegShrunk(p.Rn) {
		out = append(out, NewMulParams(p.Rd, r, p.Rm))
	}

	for _, r := range a64.RegShrunk(p.Rm) {
		out = append(out, NewMulParams(p.Rd, p.Rn, r))
	}

	return slices.Values(out)
}
