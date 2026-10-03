package arm64

// Generator for fcvtzs — one constructor (Fcvtzs) over the shared
// conversion core (rd is the gpr, rn the FP register).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// FcvtzsParams — parameters of fcvtzs.
type FcvtzsParams struct {
	FGprParams
}

func NewFcvtzsParams(p FGprParams) FcvtzsParams {
	return FcvtzsParams{FGprParams: p}
}

func (p FcvtzsParams) Instr() arm64.Instr {
	in, err := arm64.New().Fcvtzs(p.R, p.F)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p FcvtzsParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Fcvtzs — an arbitrary fcvtzs.
func Fcvtzs(rnd *rand.Rand) ohsnap.Arbitrary[FcvtzsParams] {
	base := fgprAny(rnd)
	return fcvtzsArb{base: base}
}

type fcvtzsArb struct {
	base fgprGen
}

func (a fcvtzsArb) Generate() iter.Seq[FcvtzsParams] {
	return arbStream(func() FcvtzsParams {
		return NewFcvtzsParams(ohsnap.First(a.base.Generate()))
	})
}

func (a fcvtzsArb) Shrink(p FcvtzsParams) iter.Seq[FcvtzsParams] {
	shrinks := slices.Collect(a.base.Shrink(p.FGprParams))
	out := make([]FcvtzsParams, 0, len(shrinks))
	for _, s := range shrinks {
		out = append(out, NewFcvtzsParams(s))
	}

	return slices.Values(out)
}
