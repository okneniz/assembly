package arm64

// Generator for cmtst — one constructor (Cmtst) over the shared
// three-register vector core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// CmtstParams — parameters of cmtst.
type CmtstParams struct {
	V3Params
}

func NewCmtstParams(p V3Params) CmtstParams {
	return CmtstParams{V3Params: p}
}

func (p CmtstParams) Instr() arm64.Instr {
	in, err := arm64.New().Cmtst(p.Rd, p.Rn, p.Rm, p.Arr)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p CmtstParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Cmtst — an arbitrary cmtst.
func Cmtst(rnd *rand.Rand) ohsnap.Arbitrary[CmtstParams] {
	base := newV3Gen(rnd, arrFull())
	return cmtstArb{base: base}
}

type cmtstArb struct {
	base v3Gen
}

func (a cmtstArb) Generate() iter.Seq[CmtstParams] {
	return arbStream(func() CmtstParams {
		return NewCmtstParams(ohsnap.First(a.base.Generate()))
	})
}

func (a cmtstArb) Shrink(p CmtstParams) iter.Seq[CmtstParams] {
	var out []CmtstParams
	for _, s := range slices.Collect(a.base.Shrink(p.V3Params)) {
		out = append(out, NewCmtstParams(s))
	}

	return slices.Values(out)
}
