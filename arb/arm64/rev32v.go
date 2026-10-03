package arm64

// Generator for rev32v — one constructor (Rev32V) over the shared
// two-register vector core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// Rev32VParams — parameters of rev32v.
type Rev32VParams struct {
	V2Params
}

func NewRev32VParams(p V2Params) Rev32VParams {
	return Rev32VParams{V2Params: p}
}

func (p Rev32VParams) Instr() arm64.Instr {
	in, err := arm64.New().Rev32V(p.Rd, p.Rn, p.Arr)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p Rev32VParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Rev32V — an arbitrary rev32v.
func Rev32V(rnd *rand.Rand) ohsnap.Arbitrary[Rev32VParams] {
	base := newV2Gen(rnd, arrHalf())
	return rev32vArb{base: base}
}

type rev32vArb struct {
	base v2Gen
}

func (a rev32vArb) Generate() iter.Seq[Rev32VParams] {
	return arbStream(func() Rev32VParams {
		return NewRev32VParams(ohsnap.First(a.base.Generate()))
	})
}

func (a rev32vArb) Shrink(p Rev32VParams) iter.Seq[Rev32VParams] {
	shrinks := slices.Collect(a.base.Shrink(p.V2Params))
	out := make([]Rev32VParams, 0, len(shrinks))
	for _, s := range shrinks {
		out = append(out, NewRev32VParams(s))
	}

	return slices.Values(out)
}
