package arm64

// Generator for ubfm — a thin unsigned variant over the Bfm family: one
// constructor (Ubfm), the generation and the shrink of the base bfm.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// UbfmParams — parameters of ubfm rd, rn, #immr, #imms.
type UbfmParams struct {
	BfmParams
}

func NewUbfmParams(p BfmParams) UbfmParams {
	return UbfmParams{BfmParams: p}
}

func (p UbfmParams) Instr() arm64.Instr {
	in, err := arm64.New().Ubfm(p.Rd, p.Rn, p.Immr, p.Imms)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p UbfmParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Ubfm — an arbitrary ubfm.
func Ubfm(rnd *rand.Rand) ohsnap.Arbitrary[UbfmParams] {
	base := Bfm(rnd)
	return ubfmArb{base: base}
}

type ubfmArb struct {
	base ohsnap.Arbitrary[BfmParams]
}

func (a ubfmArb) Generate() iter.Seq[UbfmParams] {
	return arbStream(func() UbfmParams {
		return NewUbfmParams(ohsnap.First(a.base.Generate()))
	})
}

func (a ubfmArb) Shrink(p UbfmParams) iter.Seq[UbfmParams] {
	shrinks := slices.Collect(a.base.Shrink(p.BfmParams))
	out := make([]UbfmParams, 0, len(shrinks))
	for _, s := range shrinks {
		out = append(out, NewUbfmParams(s))
	}

	return slices.Values(out)
}
