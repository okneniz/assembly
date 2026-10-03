package alias

// Generator for the cmn alias — one generator, one type, one text form
// family: cmn rn, #imm12[, lsl #12] | cmn rn, rm.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	a64 "github.com/okneniz/assembly/arb/arm64"
	arm64 "github.com/okneniz/assembly/arch/arm64"
)

// CmnParams — parameters of the cmn alias.
type CmnParams struct {
	Rn, Rm arm64.Reg
	Imm    int64
	Sh     bool // lsl #12
	IsImm  bool
}

func NewCmnParams(rn arm64.Reg, rm arm64.Reg, imm int64, sh bool, isImm bool) CmnParams {
	return CmnParams{
		Rn:    rn,
		Rm:    rm,
		Imm:   imm,
		Sh:    sh,
		IsImm: isImm,
	}
}

func (p CmnParams) Instr() arm64.Instr {
	in, err := instrOfText(p.String())
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

func (p CmnParams) String() string {
	if !p.IsImm {
		return "cmn " + p.Rn.String() + ", " + p.Rm.String()
	}

	if p.Sh {
		return "cmn " + p.Rn.String() + ", #" + itoa(p.Imm) + ", lsl #12"
	}

	return "cmn " + p.Rn.String() + ", #" + itoa(p.Imm)
}

// cmnGen — generator for cmn: same-width registers, immediate 0..4095.
type cmnGen struct {
	rnd *rand.Rand
}

// Cmn — an arbitrary cmn.
func Cmn(rnd *rand.Rand) ohsnap.Arbitrary[CmnParams] {
	return newCmnGen(rnd)
}

func newCmnGen(rnd *rand.Rand) cmnGen {
	return cmnGen{rnd: rnd}
}

func (g cmnGen) Generate() iter.Seq[CmnParams] {
	return stream(func() CmnParams {
		is64 := g.rnd.IntN(2) == 1
		return NewCmnParams(
			a64.GenReg(g.rnd, is64, false, true),
			a64.GenReg(g.rnd, is64, false, true),
			g.rnd.Int64N(0x1000),
			g.rnd.IntN(2) == 1,
			g.rnd.IntN(2) == 1,
		)
	})
}

func (g cmnGen) Shrink(p CmnParams) iter.Seq[CmnParams] {
	var out []CmnParams
	for _, r := range a64.RegShrunk(p.Rn) {
		out = append(out, NewCmnParams(r, p.Rm, p.Imm, p.Sh, p.IsImm))
	}

	for _, r := range a64.RegShrunk(p.Rm) {
		out = append(out, NewCmnParams(p.Rn, r, p.Imm, p.Sh, p.IsImm))
	}

	for _, v := range halved(p.Imm) {
		out = append(out, NewCmnParams(p.Rn, p.Rm, v, p.Sh, p.IsImm))
	}

	if p.Sh {
		out = append(out, NewCmnParams(p.Rn, p.Rm, p.Imm, false, true))
	}

	if p.IsImm {
		out = append(out, NewCmnParams(p.Rn, p.Rn, 0, false, false))
	}

	return slices.Values(out)
}
