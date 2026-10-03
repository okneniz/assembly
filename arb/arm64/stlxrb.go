package arm64

// Generator for stlxrb — one constructor (Stlxrb) over the shared
// exclusive core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// StlxrbParams — parameters of the stlxrb form.
type StlxrbParams struct {
	ExclParams
}

func NewStlxrbParams(p ExclParams) StlxrbParams {
	return StlxrbParams{ExclParams: p}
}

func (p StlxrbParams) Instr() arm64.Instr {
	in, err := arm64.New().Stlxrb(p.Rs, p.Rt, p.Rn)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p StlxrbParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Stlxrb — an arbitrary stlxrb.
func Stlxrb(rnd *rand.Rand) ohsnap.Arbitrary[StlxrbParams] {
	base := exclW(rnd)
	return stlxrbArb{base: base}
}

type stlxrbArb struct {
	base exclGen
}

func (a stlxrbArb) Generate() iter.Seq[StlxrbParams] {
	return arbStream(func() StlxrbParams {
		return NewStlxrbParams(ohsnap.First(a.base.Generate()))
	})
}

func (a stlxrbArb) Shrink(p StlxrbParams) iter.Seq[StlxrbParams] {
	shrinks := slices.Collect(a.base.Shrink(p.ExclParams))
	out := make([]StlxrbParams, 0, len(shrinks))
	for _, s := range shrinks {
		out = append(out, NewStlxrbParams(s))
	}

	return slices.Values(out)
}
