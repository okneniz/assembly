package arm64

// Generator for abs — one constructor (Abs) over the shared
// two-register vector core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// AbsParams — parameters of abs.
type AbsParams struct {
	V2Params
}

func NewAbsParams(p V2Params) AbsParams {
	return AbsParams{V2Params: p}
}

func (p AbsParams) Instr() arm64.Instr {
	in, err := arm64.New().Abs(p.Rd, p.Rn, p.Arr)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p AbsParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Abs — an arbitrary abs.
func Abs(rnd *rand.Rand) ohsnap.Arbitrary[AbsParams] {
	base := newV2Gen(rnd, arrFull())
	return absArb{base: base}
}

type absArb struct {
	base v2Gen
}

func (a absArb) Generate() iter.Seq[AbsParams] {
	return arbStream(func() AbsParams {
		return NewAbsParams(ohsnap.First(a.base.Generate()))
	})
}

func (a absArb) Shrink(p AbsParams) iter.Seq[AbsParams] {
	var out []AbsParams
	for _, s := range slices.Collect(a.base.Shrink(p.V2Params)) {
		out = append(out, NewAbsParams(s))
	}

	return slices.Values(out)
}
