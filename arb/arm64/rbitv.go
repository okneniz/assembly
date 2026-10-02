package arm64

// Generator for rbitv — one constructor (RbitV) over the shared
// two-register vector core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// RbitVParams — parameters of rbitv.
type RbitVParams struct {
	V2Params
}

func NewRbitVParams(p V2Params) RbitVParams {
	return RbitVParams{V2Params: p}
}

func (p RbitVParams) Instr() arm64.Instr {
	in, err := arm64.New().RbitV(p.Rd, p.Rn, p.Arr)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p RbitVParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// RbitV — an arbitrary rbitv.
func RbitV(rnd *rand.Rand) ohsnap.Arbitrary[RbitVParams] {
	base := newV2Gen(rnd, arrLogical())
	return rbitvArb{base: base}
}

type rbitvArb struct {
	base v2Gen
}

func (a rbitvArb) Generate() iter.Seq[RbitVParams] {
	return arbStream(func() RbitVParams {
		return NewRbitVParams(ohsnap.First(a.base.Generate()))
	})
}

func (a rbitvArb) Shrink(p RbitVParams) iter.Seq[RbitVParams] {
	var out []RbitVParams
	for _, s := range slices.Collect(a.base.Shrink(p.V2Params)) {
		out = append(out, NewRbitVParams(s))
	}

	return slices.Values(out)
}
