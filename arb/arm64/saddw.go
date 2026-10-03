package arm64

// Generator for saddw — one constructor (Saddw) over the shared
// three-register vector core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// SaddwParams — parameters of saddw.
type SaddwParams struct {
	V3Params
}

func NewSaddwParams(p V3Params) SaddwParams {
	return SaddwParams{V3Params: p}
}

func (p SaddwParams) Instr() arm64.Instr {
	in, err := arm64.New().Saddw(p.Rd, p.Rn, p.Rm, p.Arr)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p SaddwParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Saddw — an arbitrary saddw.
func Saddw(rnd *rand.Rand) ohsnap.Arbitrary[SaddwParams] {
	base := newV3Gen(rnd, arrWiden())
	return saddwArb{base: base}
}

type saddwArb struct {
	base v3Gen
}

func (a saddwArb) Generate() iter.Seq[SaddwParams] {
	return arbStream(func() SaddwParams {
		return NewSaddwParams(ohsnap.First(a.base.Generate()))
	})
}

func (a saddwArb) Shrink(p SaddwParams) iter.Seq[SaddwParams] {
	shrinks := slices.Collect(a.base.Shrink(p.V3Params))
	out := make([]SaddwParams, 0, len(shrinks))
	for _, s := range shrinks {
		out = append(out, NewSaddwParams(s))
	}

	return slices.Values(out)
}
