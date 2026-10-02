package arm64

// Generator for stxrb — one constructor (Stxrb) over the shared
// exclusive core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// StxrbParams — parameters of the stxrb form.
type StxrbParams struct {
	ExclParams
}

func NewStxrbParams(p ExclParams) StxrbParams {
	return StxrbParams{ExclParams: p}
}

func (p StxrbParams) Instr() arm64.Instr {
	in, err := arm64.New().Stxrb(p.Rs, p.Rt, p.Rn)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p StxrbParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Stxrb — an arbitrary stxrb (the b form: rt is a w register).
func Stxrb(rnd *rand.Rand) ohsnap.Arbitrary[StxrbParams] {
	base := exclW(rnd)
	return stxrbArb{base: base}
}

type stxrbArb struct {
	base exclGen
}

func (a stxrbArb) Generate() iter.Seq[StxrbParams] {
	return arbStream(func() StxrbParams {
		return NewStxrbParams(ohsnap.First(a.base.Generate()))
	})
}

func (a stxrbArb) Shrink(p StxrbParams) iter.Seq[StxrbParams] {
	var out []StxrbParams
	for _, s := range slices.Collect(a.base.Shrink(p.ExclParams)) {
		out = append(out, NewStxrbParams(s))
	}

	return slices.Values(out)
}
