package arm64

// Generator for fmadd — one constructor (Fmadd) over the shared
// four-register FP core (one s or d kind).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// FmaddParams — parameters of fmadd.
type FmaddParams struct {
	F4Params
}

func NewFmaddParams(p F4Params) FmaddParams {
	return FmaddParams{F4Params: p}
}

func (p FmaddParams) Instr() arm64.Instr {
	in, err := arm64.New().Fmadd(p.Rd, p.Rn, p.Rm, p.Ra)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p FmaddParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Fmadd — an arbitrary fmadd.
func Fmadd(rnd *rand.Rand) ohsnap.Arbitrary[FmaddParams] {
	base := f4(rnd)
	return fmaddArb{base: base}
}

type fmaddArb struct {
	base f4Gen
}

func (a fmaddArb) Generate() iter.Seq[FmaddParams] {
	return arbStream(func() FmaddParams {
		return NewFmaddParams(ohsnap.First(a.base.Generate()))
	})
}

func (a fmaddArb) Shrink(p FmaddParams) iter.Seq[FmaddParams] {
	shrinks := slices.Collect(a.base.Shrink(p.F4Params))
	out := make([]FmaddParams, 0, len(shrinks))
	for _, s := range shrinks {
		out = append(out, NewFmaddParams(s))
	}

	return slices.Values(out)
}
