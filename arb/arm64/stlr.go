package arm64

// Generator for stlr — one constructor (Stlr) over the shared
// atomic core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// StlrParams — parameters of the stlr form.
type StlrParams struct {
	LsParams
}

func NewStlrParams(p LsParams) StlrParams {
	return StlrParams{LsParams: p}
}

func (p StlrParams) Instr() arm64.Instr {
	in, err := arm64.New().Stlr(p.Rt, p.Rn)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p StlrParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Stlr — an arbitrary stlr.
func Stlr(rnd *rand.Rand) ohsnap.Arbitrary[StlrParams] {
	base := atom(rnd)
	return stlrArb{base: base}
}

type stlrArb struct {
	base atomGen
}

func (a stlrArb) Generate() iter.Seq[StlrParams] {
	return arbStream(func() StlrParams {
		return NewStlrParams(ohsnap.First(a.base.Generate()))
	})
}

func (a stlrArb) Shrink(p StlrParams) iter.Seq[StlrParams] {
	shrinks := slices.Collect(a.base.Shrink(p.LsParams))
	out := make([]StlrParams, 0, len(shrinks))
	for _, s := range shrinks {
		out = append(out, NewStlrParams(s))
	}

	return slices.Values(out)
}
