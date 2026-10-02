package arm64

// Generator for bic — one constructor (Bic) over the shared
// three-register vector core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// BicVParams — parameters of bic.
type BicVParams struct {
	V3Params
}

func NewBicVParams(p V3Params) BicVParams {
	return BicVParams{V3Params: p}
}

func (p BicVParams) Instr() arm64.Instr {
	in, err := arm64.New().Bic(p.Rd, p.Rn, p.Rm, p.Arr)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p BicVParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Bic — an arbitrary bic.
func Bic(rnd *rand.Rand) ohsnap.Arbitrary[BicVParams] {
	base := newV3Gen(rnd, arrLogical())
	return bicArb{base: base}
}

type bicArb struct {
	base v3Gen
}

func (a bicArb) Generate() iter.Seq[BicVParams] {
	return arbStream(func() BicVParams {
		return NewBicVParams(ohsnap.First(a.base.Generate()))
	})
}

func (a bicArb) Shrink(p BicVParams) iter.Seq[BicVParams] {
	var out []BicVParams
	for _, s := range slices.Collect(a.base.Shrink(p.V3Params)) {
		out = append(out, NewBicVParams(s))
	}

	return slices.Values(out)
}
