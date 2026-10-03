package arm64

// Generator for fmax — one constructor (Fmax) over the shared
// three-register FP core (one s or d kind).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// FmaxParams — parameters of fmax.
type FmaxParams struct {
	F3Params
}

func NewFmaxParams(p F3Params) FmaxParams {
	return FmaxParams{F3Params: p}
}

func (p FmaxParams) Instr() arm64.Instr {
	in, err := arm64.New().Fmax(p.Rd, p.Rn, p.Rm)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p FmaxParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Fmax — an arbitrary fmax.
func Fmax(rnd *rand.Rand) ohsnap.Arbitrary[FmaxParams] {
	base := f3(rnd)
	return fmaxArb{base: base}
}

type fmaxArb struct {
	base f3Gen
}

func (a fmaxArb) Generate() iter.Seq[FmaxParams] {
	return arbStream(func() FmaxParams {
		return NewFmaxParams(ohsnap.First(a.base.Generate()))
	})
}

func (a fmaxArb) Shrink(p FmaxParams) iter.Seq[FmaxParams] {
	shrinks := slices.Collect(a.base.Shrink(p.F3Params))
	out := make([]FmaxParams, 0, len(shrinks))
	for _, s := range shrinks {
		out = append(out, NewFmaxParams(s))
	}

	return slices.Values(out)
}
