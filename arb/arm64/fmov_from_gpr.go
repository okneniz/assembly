package arm64

// Generator for fmov_from_gpr — one constructor (FmovFromGpr) over the shared
// conversion core (rd is the FP register, rn the gpr).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// FmovFromGprParams — parameters of fmov_from_gpr.
type FmovFromGprParams struct {
	FGprParams
}

func NewFmovFromGprParams(p FGprParams) FmovFromGprParams {
	return FmovFromGprParams{FGprParams: p}
}

func (p FmovFromGprParams) Instr() arm64.Instr {
	in, err := arm64.New().FmovFromGpr(p.F, p.R)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p FmovFromGprParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// FmovFromGpr — an arbitrary fmov_from_gpr.
func FmovFromGpr(rnd *rand.Rand) ohsnap.Arbitrary[FmovFromGprParams] {
	base := fgpr(rnd)
	return fmov_from_gprArb{base: base}
}

type fmov_from_gprArb struct {
	base fgprGen
}

func (a fmov_from_gprArb) Generate() iter.Seq[FmovFromGprParams] {
	return arbStream(func() FmovFromGprParams {
		return NewFmovFromGprParams(ohsnap.First(a.base.Generate()))
	})
}

func (a fmov_from_gprArb) Shrink(p FmovFromGprParams) iter.Seq[FmovFromGprParams] {
	shrinks := slices.Collect(a.base.Shrink(p.FGprParams))
	out := make([]FmovFromGprParams, 0, len(shrinks))
	for _, s := range shrinks {
		out = append(out, NewFmovFromGprParams(s))
	}

	return slices.Values(out)
}
