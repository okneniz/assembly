package arm64

// Generator for prfm — one generator, one type, one constructor (Prfm):
// a single operand, the base register (the hint is fixed by the ctor).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// PrfmParams — parameters of the prfm form.
type PrfmParams struct {
	Rn arm64.Reg
}

func NewPrfmParams(rn arm64.Reg) PrfmParams {
	return PrfmParams{Rn: rn}
}

func (p PrfmParams) Instr() arm64.Instr {
	in, err := arm64.New().Prfm(p.Rn)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

func (p PrfmParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// prfmGen — generator for prfm: an x/sp base.
type prfmGen struct {
	rnd *rand.Rand
}

// Prfm — an arbitrary prfm.
func Prfm(rnd *rand.Rand) ohsnap.Arbitrary[PrfmParams] {
	return newPrfmGen(rnd)
}

func newPrfmGen(rnd *rand.Rand) prfmGen {
	return prfmGen{rnd: rnd}
}

func (g prfmGen) Generate() iter.Seq[PrfmParams] {
	return arbStream(func() PrfmParams {
		return NewPrfmParams(genReg(g.rnd, true, true, false))
	})
}

func (g prfmGen) Shrink(p PrfmParams) iter.Seq[PrfmParams] {
	regs := regShrunk(p.Rn)
	out := make([]PrfmParams, 0, len(regs))
	for _, r := range regs {
		out = append(out, NewPrfmParams(r))
	}

	return slices.Values(out)
}
