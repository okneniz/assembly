package arm64

// Generator for ssubw — one constructor (Ssubw) over the shared
// three-register vector core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// SsubwParams — parameters of ssubw.
type SsubwParams struct {
	V3Params
}

func NewSsubwParams(p V3Params) SsubwParams {
	return SsubwParams{V3Params: p}
}

func (p SsubwParams) Instr() arm64.Instr {
	in, err := arm64.New().Ssubw(p.Rd, p.Rn, p.Rm, p.Arr)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p SsubwParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Ssubw — an arbitrary ssubw.
func Ssubw(rnd *rand.Rand) ohsnap.Arbitrary[SsubwParams] {
	base := newV3Gen(rnd, arrWiden())
	return ssubwArb{base: base}
}

type ssubwArb struct {
	base v3Gen
}

func (a ssubwArb) Generate() iter.Seq[SsubwParams] {
	return arbStream(func() SsubwParams {
		return NewSsubwParams(ohsnap.First(a.base.Generate()))
	})
}

func (a ssubwArb) Shrink(p SsubwParams) iter.Seq[SsubwParams] {
	var out []SsubwParams
	for _, s := range slices.Collect(a.base.Shrink(p.V3Params)) {
		out = append(out, NewSsubwParams(s))
	}

	return slices.Values(out)
}
