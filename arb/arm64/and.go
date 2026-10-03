package arm64

// Generator for and — one constructor (And) over the shared
// three-register vector core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// AndVParams — parameters of and.
type AndVParams struct {
	V3Params
}

func NewAndVParams(p V3Params) AndVParams {
	return AndVParams{V3Params: p}
}

func (p AndVParams) Instr() arm64.Instr {
	in, err := arm64.New().And(p.Rd, p.Rn, p.Rm, p.Arr)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p AndVParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// And — an arbitrary and.
func And(rnd *rand.Rand) ohsnap.Arbitrary[AndVParams] {
	base := newV3Gen(rnd, arrLogical())
	return andArb{base: base}
}

type andArb struct {
	base v3Gen
}

func (a andArb) Generate() iter.Seq[AndVParams] {
	return arbStream(func() AndVParams {
		return NewAndVParams(ohsnap.First(a.base.Generate()))
	})
}

func (a andArb) Shrink(p AndVParams) iter.Seq[AndVParams] {
	shrinks := slices.Collect(a.base.Shrink(p.V3Params))
	out := make([]AndVParams, 0, len(shrinks))
	for _, s := range shrinks {
		out = append(out, NewAndVParams(s))
	}

	return slices.Values(out)
}
