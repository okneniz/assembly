package arm64

// Generator for fmov_to_gpr — one constructor (FmovToGpr) over the shared
// conversion core (rd is the gpr, rn the FP register).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// FmovToGprParams — parameters of fmov_to_gpr.
type FmovToGprParams struct {
	FGprParams
}

func NewFmovToGprParams(p FGprParams) FmovToGprParams {
	return FmovToGprParams{FGprParams: p}
}

func (p FmovToGprParams) Instr() arm64.Instr {
	in, err := arm64.New().FmovToGpr(p.R, p.F)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p FmovToGprParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// FmovToGpr — an arbitrary fmov_to_gpr.
func FmovToGpr(rnd *rand.Rand) ohsnap.Arbitrary[FmovToGprParams] {
	base := fgpr(rnd)
	return fmov_to_gprArb{base: base}
}

type fmov_to_gprArb struct {
	base fgprGen
}

func (a fmov_to_gprArb) Generate() iter.Seq[FmovToGprParams] {
	return arbStream(func() FmovToGprParams {
		return NewFmovToGprParams(ohsnap.First(a.base.Generate()))
	})
}

func (a fmov_to_gprArb) Shrink(p FmovToGprParams) iter.Seq[FmovToGprParams] {
	shrinks := slices.Collect(a.base.Shrink(p.FGprParams))
	out := make([]FmovToGprParams, 0, len(shrinks))
	for _, s := range shrinks {
		out = append(out, NewFmovToGprParams(s))
	}

	return slices.Values(out)
}
