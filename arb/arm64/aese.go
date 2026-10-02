package arm64

// Generator for aese — one constructor (Aese) over the shared
// arrangement-less vector core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// AeseParams — parameters of aese.
type AeseParams struct {
	V2PlainParams
}

func NewAeseParams(p V2PlainParams) AeseParams {
	return AeseParams{V2PlainParams: p}
}

func (p AeseParams) Instr() arm64.Instr {
	in, err := arm64.New().Aese(p.Rd, p.Rn)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p AeseParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Aese — an arbitrary aese.
func Aese(rnd *rand.Rand) ohsnap.Arbitrary[AeseParams] {
	base := v2Plain(rnd)
	return aeseArb{base: base}
}

type aeseArb struct {
	base v2PlainGen
}

func (a aeseArb) Generate() iter.Seq[AeseParams] {
	return arbStream(func() AeseParams {
		return NewAeseParams(ohsnap.First(a.base.Generate()))
	})
}

func (a aeseArb) Shrink(p AeseParams) iter.Seq[AeseParams] {
	var out []AeseParams
	for _, s := range slices.Collect(a.base.Shrink(p.V2PlainParams)) {
		out = append(out, NewAeseParams(s))
	}

	return slices.Values(out)
}
