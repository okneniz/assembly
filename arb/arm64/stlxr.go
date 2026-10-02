package arm64

// Generator for stlxr — one constructor (Stlxr) over the shared
// exclusive core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// StlxrParams — parameters of the stlxr form.
type StlxrParams struct {
	ExclParams
}

func NewStlxrParams(p ExclParams) StlxrParams {
	return StlxrParams{ExclParams: p}
}

func (p StlxrParams) Instr() arm64.Instr {
	in, err := arm64.New().Stlxr(p.Rs, p.Rt, p.Rn)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p StlxrParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Stlxr — an arbitrary stlxr.
func Stlxr(rnd *rand.Rand) ohsnap.Arbitrary[StlxrParams] {
	base := excl(rnd)
	return stlxrArb{base: base}
}

type stlxrArb struct {
	base exclGen
}

func (a stlxrArb) Generate() iter.Seq[StlxrParams] {
	return arbStream(func() StlxrParams {
		return NewStlxrParams(ohsnap.First(a.base.Generate()))
	})
}

func (a stlxrArb) Shrink(p StlxrParams) iter.Seq[StlxrParams] {
	var out []StlxrParams
	for _, s := range slices.Collect(a.base.Shrink(p.ExclParams)) {
		out = append(out, NewStlxrParams(s))
	}

	return slices.Values(out)
}
