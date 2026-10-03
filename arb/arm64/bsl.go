package arm64

// Generator for bsl — one constructor (Bsl) over the shared
// three-register vector core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// BslParams — parameters of bsl.
type BslParams struct {
	V3Params
}

func NewBslParams(p V3Params) BslParams {
	return BslParams{V3Params: p}
}

func (p BslParams) Instr() arm64.Instr {
	in, err := arm64.New().Bsl(p.Rd, p.Rn, p.Rm, p.Arr)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p BslParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Bsl — an arbitrary bsl.
func Bsl(rnd *rand.Rand) ohsnap.Arbitrary[BslParams] {
	base := newV3Gen(rnd, arrLogical())
	return bslArb{base: base}
}

type bslArb struct {
	base v3Gen
}

func (a bslArb) Generate() iter.Seq[BslParams] {
	return arbStream(func() BslParams {
		return NewBslParams(ohsnap.First(a.base.Generate()))
	})
}

func (a bslArb) Shrink(p BslParams) iter.Seq[BslParams] {
	shrinks := slices.Collect(a.base.Shrink(p.V3Params))
	out := make([]BslParams, 0, len(shrinks))
	for _, s := range shrinks {
		out = append(out, NewBslParams(s))
	}

	return slices.Values(out)
}
