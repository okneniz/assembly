package arm64

// Generator for dup — one constructor (Dup): the gpr source fans out
// over the lanes (w for the b/h/s lanes, x for d).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// DupGenParams — parameters of dup.
type DupGenParams struct {
	DupParams
}

func NewDupGenParams(p DupParams) DupGenParams {
	return DupGenParams{DupParams: p}
}

func (p DupGenParams) Instr() arm64.Instr {
	in, err := arm64.New().Dup(p.Rd, p.Wn, p.Arr)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p DupGenParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// DupGen — an arbitrary dup (the generator family of the gpr source).
func DupGen(rnd *rand.Rand) ohsnap.Arbitrary[DupGenParams] {
	base := newDupGen(rnd)
	return dupArb{base: base}
}

type dupArb struct {
	base dupGen
}

func (a dupArb) Generate() iter.Seq[DupGenParams] {
	return arbStream(func() DupGenParams {
		return NewDupGenParams(ohsnap.First(a.base.Generate()))
	})
}

func (a dupArb) Shrink(p DupGenParams) iter.Seq[DupGenParams] {
	shrinks := slices.Collect(a.base.Shrink(p.DupParams))
	out := make([]DupGenParams, 0, len(shrinks))
	for _, s := range shrinks {
		out = append(out, NewDupGenParams(s))
	}

	return slices.Values(out)
}
