package arm64

// Generator for aesmc — one constructor (Aesmc) over the shared
// arrangement-less vector core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// AesmcParams — parameters of aesmc.
type AesmcParams struct {
	V2PlainParams
}

func NewAesmcParams(p V2PlainParams) AesmcParams {
	return AesmcParams{V2PlainParams: p}
}

func (p AesmcParams) Instr() arm64.Instr {
	in, err := arm64.New().Aesmc(p.Rd, p.Rn)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p AesmcParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Aesmc — an arbitrary aesmc.
func Aesmc(rnd *rand.Rand) ohsnap.Arbitrary[AesmcParams] {
	base := v2Plain(rnd)
	return aesmcArb{base: base}
}

type aesmcArb struct {
	base v2PlainGen
}

func (a aesmcArb) Generate() iter.Seq[AesmcParams] {
	return arbStream(func() AesmcParams {
		return NewAesmcParams(ohsnap.First(a.base.Generate()))
	})
}

func (a aesmcArb) Shrink(p AesmcParams) iter.Seq[AesmcParams] {
	var out []AesmcParams
	for _, s := range slices.Collect(a.base.Shrink(p.V2PlainParams)) {
		out = append(out, NewAesmcParams(s))
	}

	return slices.Values(out)
}
