package arm64

// Generator for stur — one constructor (Stur) over the shared
// load/store core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// SturParams — parameters of the stur form.
type SturParams struct {
	LsParams
}

func NewSturParams(p LsParams) SturParams {
	return SturParams{LsParams: p}
}

func (p SturParams) Instr() arm64.Instr {
	in, err := arm64.New().Stur(p.Rt, p.Rn, p.Off)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p SturParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Stur — an arbitrary stur.
func Stur(rnd *rand.Rand) ohsnap.Arbitrary[SturParams] {
	base := newLsGen(rnd, nil, UnscaledOff(rnd))
	return sturArb{base: base}
}

type sturArb struct {
	base lsGen
}

func (a sturArb) Generate() iter.Seq[SturParams] {
	return arbStream(func() SturParams {
		return NewSturParams(ohsnap.First(a.base.Generate()))
	})
}

func (a sturArb) Shrink(p SturParams) iter.Seq[SturParams] {
	shrinks := slices.Collect(a.base.Shrink(p.LsParams))
	out := make([]SturParams, 0, len(shrinks))
	for _, s := range shrinks {
		out = append(out, NewSturParams(s))
	}

	return slices.Values(out)
}
