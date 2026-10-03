package arm64

// Generator for cmge — one constructor (Cmge) over the shared
// three-register vector core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// CmgeParams — parameters of cmge.
type CmgeParams struct {
	V3Params
}

func NewCmgeParams(p V3Params) CmgeParams {
	return CmgeParams{V3Params: p}
}

func (p CmgeParams) Instr() arm64.Instr {
	in, err := arm64.New().Cmge(p.Rd, p.Rn, p.Rm, p.Arr)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p CmgeParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Cmge — an arbitrary cmge.
func Cmge(rnd *rand.Rand) ohsnap.Arbitrary[CmgeParams] {
	base := newV3Gen(rnd, arrFull())
	return cmgeArb{base: base}
}

type cmgeArb struct {
	base v3Gen
}

func (a cmgeArb) Generate() iter.Seq[CmgeParams] {
	return arbStream(func() CmgeParams {
		return NewCmgeParams(ohsnap.First(a.base.Generate()))
	})
}

func (a cmgeArb) Shrink(p CmgeParams) iter.Seq[CmgeParams] {
	shrinks := slices.Collect(a.base.Shrink(p.V3Params))
	out := make([]CmgeParams, 0, len(shrinks))
	for _, s := range shrinks {
		out = append(out, NewCmgeParams(s))
	}

	return slices.Values(out)
}
