package arm64

// Generator for sbfm — a thin signed variant over the Bfm family: one
// constructor (Sbfm), the generation and the shrink of the base bfm.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// SbfmParams — parameters of sbfm rd, rn, #immr, #imms.
type SbfmParams struct {
	BfmParams
}

func NewSbfmParams(p BfmParams) SbfmParams {
	return SbfmParams{BfmParams: p}
}

func (p SbfmParams) Instr() arm64.Instr {
	in, err := arm64.New().Sbfm(p.Rd, p.Rn, p.Immr, p.Imms)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p SbfmParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Sbfm — an arbitrary sbfm.
func Sbfm(rnd *rand.Rand) ohsnap.Arbitrary[SbfmParams] {
	base := Bfm(rnd)
	return sbfmArb{base: base}
}

type sbfmArb struct {
	base ohsnap.Arbitrary[BfmParams]
}

func (a sbfmArb) Generate() iter.Seq[SbfmParams] {
	return arbStream(func() SbfmParams {
		return NewSbfmParams(ohsnap.First(a.base.Generate()))
	})
}

func (a sbfmArb) Shrink(p SbfmParams) iter.Seq[SbfmParams] {
	shrinks := slices.Collect(a.base.Shrink(p.BfmParams))
	out := make([]SbfmParams, 0, len(shrinks))
	for _, s := range shrinks {
		out = append(out, NewSbfmParams(s))
	}

	return slices.Values(out)
}
