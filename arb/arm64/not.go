package arm64

// Generator for not — one constructor (Not) over the shared
// two-register vector core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// NotParams — parameters of not.
type NotParams struct {
	V2Params
}

func NewNotParams(p V2Params) NotParams {
	return NotParams{V2Params: p}
}

func (p NotParams) Instr() arm64.Instr {
	in, err := arm64.New().Not(p.Rd, p.Rn, p.Arr)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p NotParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Not — an arbitrary not.
func Not(rnd *rand.Rand) ohsnap.Arbitrary[NotParams] {
	base := newV2Gen(rnd, arrLogical())
	return notArb{base: base}
}

type notArb struct {
	base v2Gen
}

func (a notArb) Generate() iter.Seq[NotParams] {
	return arbStream(func() NotParams {
		return NewNotParams(ohsnap.First(a.base.Generate()))
	})
}

func (a notArb) Shrink(p NotParams) iter.Seq[NotParams] {
	var out []NotParams
	for _, s := range slices.Collect(a.base.Shrink(p.V2Params)) {
		out = append(out, NewNotParams(s))
	}

	return slices.Values(out)
}
