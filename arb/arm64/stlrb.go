package arm64

// Generator for stlrb — one constructor (Stlrb) over the shared
// atomic core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// StlrbParams — parameters of the stlrb form.
type StlrbParams struct {
	LsParams
}

func NewStlrbParams(p LsParams) StlrbParams {
	return StlrbParams{LsParams: p}
}

func (p StlrbParams) Instr() arm64.Instr {
	in, err := arm64.New().Stlrb(p.Rt, p.Rn)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p StlrbParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Stlrb — an arbitrary stlrb.
func Stlrb(rnd *rand.Rand) ohsnap.Arbitrary[StlrbParams] {
	base := atomW(rnd)
	return stlrbArb{base: base}
}

type stlrbArb struct {
	base atomGen
}

func (a stlrbArb) Generate() iter.Seq[StlrbParams] {
	return arbStream(func() StlrbParams {
		return NewStlrbParams(ohsnap.First(a.base.Generate()))
	})
}

func (a stlrbArb) Shrink(p StlrbParams) iter.Seq[StlrbParams] {
	shrinks := slices.Collect(a.base.Shrink(p.LsParams))
	out := make([]StlrbParams, 0, len(shrinks))
	for _, s := range shrinks {
		out = append(out, NewStlrbParams(s))
	}

	return slices.Values(out)
}
