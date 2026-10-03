package arm64

// Generator for orr — one constructor (Orr) over the shared
// three-register vector core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// OrrVParams — parameters of orr.
type OrrVParams struct {
	V3Params
}

func NewOrrVParams(p V3Params) OrrVParams {
	return OrrVParams{V3Params: p}
}

func (p OrrVParams) Instr() arm64.Instr {
	in, err := arm64.New().Orr(p.Rd, p.Rn, p.Rm, p.Arr)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p OrrVParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Orr — an arbitrary orr.
func Orr(rnd *rand.Rand) ohsnap.Arbitrary[OrrVParams] {
	base := newV3Gen(rnd, arrLogical())
	return orrArb{base: base}
}

type orrArb struct {
	base v3Gen
}

func (a orrArb) Generate() iter.Seq[OrrVParams] {
	return arbStream(func() OrrVParams {
		return NewOrrVParams(ohsnap.First(a.base.Generate()))
	})
}

func (a orrArb) Shrink(p OrrVParams) iter.Seq[OrrVParams] {
	shrinks := slices.Collect(a.base.Shrink(p.V3Params))
	out := make([]OrrVParams, 0, len(shrinks))
	for _, s := range shrinks {
		out = append(out, NewOrrVParams(s))
	}

	return slices.Values(out)
}
