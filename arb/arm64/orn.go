package arm64

// Generator for orn — one constructor (Orn) over the shared
// three-register vector core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// OrnVParams — parameters of orn.
type OrnVParams struct {
	V3Params
}

func NewOrnVParams(p V3Params) OrnVParams {
	return OrnVParams{V3Params: p}
}

func (p OrnVParams) Instr() arm64.Instr {
	in, err := arm64.New().Orn(p.Rd, p.Rn, p.Rm, p.Arr)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p OrnVParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Orn — an arbitrary orn.
func Orn(rnd *rand.Rand) ohsnap.Arbitrary[OrnVParams] {
	base := newV3Gen(rnd, arrLogical())
	return ornArb{base: base}
}

type ornArb struct {
	base v3Gen
}

func (a ornArb) Generate() iter.Seq[OrnVParams] {
	return arbStream(func() OrnVParams {
		return NewOrnVParams(ohsnap.First(a.base.Generate()))
	})
}

func (a ornArb) Shrink(p OrnVParams) iter.Seq[OrnVParams] {
	var out []OrnVParams
	for _, s := range slices.Collect(a.base.Shrink(p.V3Params)) {
		out = append(out, NewOrnVParams(s))
	}

	return slices.Values(out)
}
