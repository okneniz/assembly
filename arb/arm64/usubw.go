package arm64

// Generator for usubw — one constructor (Usubw) over the shared
// three-register vector core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// UsubwParams — parameters of usubw.
type UsubwParams struct {
	V3Params
}

func NewUsubwParams(p V3Params) UsubwParams {
	return UsubwParams{V3Params: p}
}

func (p UsubwParams) Instr() arm64.Instr {
	in, err := arm64.New().Usubw(p.Rd, p.Rn, p.Rm, p.Arr)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p UsubwParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Usubw — an arbitrary usubw.
func Usubw(rnd *rand.Rand) ohsnap.Arbitrary[UsubwParams] {
	base := newV3Gen(rnd, arrWiden())
	return usubwArb{base: base}
}

type usubwArb struct {
	base v3Gen
}

func (a usubwArb) Generate() iter.Seq[UsubwParams] {
	return arbStream(func() UsubwParams {
		return NewUsubwParams(ohsnap.First(a.base.Generate()))
	})
}

func (a usubwArb) Shrink(p UsubwParams) iter.Seq[UsubwParams] {
	shrinks := slices.Collect(a.base.Shrink(p.V3Params))
	out := make([]UsubwParams, 0, len(shrinks))
	for _, s := range shrinks {
		out = append(out, NewUsubwParams(s))
	}

	return slices.Values(out)
}
