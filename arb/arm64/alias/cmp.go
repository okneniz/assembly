package alias

// Generator for the cmp alias — one generator, one type, one text form
// family: cmp rn, #imm12[, lsl #12] | cmp rn, rm.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	a64 "github.com/okneniz/assembly/arb/arm64"
	arm64 "github.com/okneniz/assembly/arch/arm64"
)

// CmpParams — parameters of the cmp alias.
type CmpParams struct {
	Rn, Rm arm64.Reg
	Imm    int64
	Sh     bool // lsl #12
	IsImm  bool
}

func NewCmpParams(rn arm64.Reg, rm arm64.Reg, imm int64, sh bool, isImm bool) CmpParams {
	return CmpParams{
		Rn:    rn,
		Rm:    rm,
		Imm:   imm,
		Sh:    sh,
		IsImm: isImm,
	}
}

func (p CmpParams) String() string {
	if !p.IsImm {
		return "cmp " + p.Rn.String() + ", " + p.Rm.String()
	}

	if p.Sh {
		return "cmp " + p.Rn.String() + ", #" + itoa(p.Imm) + ", lsl #12"
	}

	return "cmp " + p.Rn.String() + ", #" + itoa(p.Imm)
}

func (p CmpParams) Instr() arm64.Instr {
	in, err := instrOfText(p.String())
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

// cmpGen — generator for cmp: same-width registers, immediate 0..4095.
type cmpGen struct {
	rnd *rand.Rand
}

func newCmpGen(rnd *rand.Rand) cmpGen {
	return cmpGen{rnd: rnd}
}

// Cmp — an arbitrary cmp.
func Cmp(rnd *rand.Rand) ohsnap.Arbitrary[CmpParams] {
	return newCmpGen(rnd)
}

func (g cmpGen) Generate() iter.Seq[CmpParams] {
	return stream(func() CmpParams {
		is64 := g.rnd.IntN(2) == 1
		return NewCmpParams(
			a64.GenReg(g.rnd, is64, false, true),
			a64.GenReg(g.rnd, is64, false, true),
			g.rnd.Int64N(0x1000),
			g.rnd.IntN(2) == 1,
			g.rnd.IntN(2) == 1,
		)
	})
}

func (g cmpGen) Shrink(p CmpParams) iter.Seq[CmpParams] {
	var out []CmpParams
	for _, r := range a64.RegShrunk(p.Rn) {
		out = append(out, NewCmpParams(r, p.Rm, p.Imm, p.Sh, p.IsImm))
	}

	for _, r := range a64.RegShrunk(p.Rm) {
		out = append(out, NewCmpParams(p.Rn, r, p.Imm, p.Sh, p.IsImm))
	}

	for _, v := range halved(p.Imm) {
		out = append(out, NewCmpParams(p.Rn, p.Rm, v, p.Sh, p.IsImm))
	}

	if p.Sh {
		out = append(out, NewCmpParams(p.Rn, p.Rm, p.Imm, false, true))
	}

	if p.IsImm {
		out = append(out, NewCmpParams(p.Rn, p.Rn, 0, false, false))
	}

	return slices.Values(out)
}
