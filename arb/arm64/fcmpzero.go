package arm64

// Generator for fcmp (the #0.0 form) — one constructor (FcmpZero).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// FcmpZeroParams — parameters of the fcmp #0.0 form.
type FcmpZeroParams struct {
	Rn arm64.FReg
}

func NewFcmpZeroParams(rn arm64.FReg) FcmpZeroParams {
	return FcmpZeroParams{Rn: rn}
}

func (p FcmpZeroParams) Instr() arm64.Instr {
	in, err := arm64.New().FcmpZero(p.Rn)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p FcmpZeroParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// fcmpZeroGen — generator for the fcmp #0.0 form.
type fcmpZeroGen struct {
	rnd *rand.Rand
}

func newFcmpZeroGen(rnd *rand.Rand) fcmpZeroGen {
	return fcmpZeroGen{rnd: rnd}
}

// FcmpZero — an arbitrary fcmp #0.0.
func FcmpZero(rnd *rand.Rand) ohsnap.Arbitrary[FcmpZeroParams] {
	return newFcmpZeroGen(rnd)
}

func (g fcmpZeroGen) Generate() iter.Seq[FcmpZeroParams] {
	return arbStream(func() FcmpZeroParams {
		return NewFcmpZeroParams(genFp(g.rnd, g.rnd.IntN(2) == 1))
	})
}

func (g fcmpZeroGen) Shrink(p FcmpZeroParams) iter.Seq[FcmpZeroParams] {
	var out []FcmpZeroParams
	for _, r := range fpShrunk(p.Rn) {
		out = append(out, NewFcmpZeroParams(r))
	}

	return slices.Values(out)
}
