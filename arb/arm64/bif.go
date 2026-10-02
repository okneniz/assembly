package arm64

// Generator for bif — one constructor (Bif) over the shared
// three-register vector core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// BifParams — parameters of bif.
type BifParams struct {
	V3Params
}

func NewBifParams(p V3Params) BifParams {
	return BifParams{V3Params: p}
}

func (p BifParams) Instr() arm64.Instr {
	in, err := arm64.New().Bif(p.Rd, p.Rn, p.Rm, p.Arr)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p BifParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Bif — an arbitrary bif.
func Bif(rnd *rand.Rand) ohsnap.Arbitrary[BifParams] {
	base := newV3Gen(rnd, arrLogical())
	return bifArb{base: base}
}

type bifArb struct {
	base v3Gen
}

func (a bifArb) Generate() iter.Seq[BifParams] {
	return arbStream(func() BifParams {
		return NewBifParams(ohsnap.First(a.base.Generate()))
	})
}

func (a bifArb) Shrink(p BifParams) iter.Seq[BifParams] {
	var out []BifParams
	for _, s := range slices.Collect(a.base.Shrink(p.V3Params)) {
		out = append(out, NewBifParams(s))
	}

	return slices.Values(out)
}
