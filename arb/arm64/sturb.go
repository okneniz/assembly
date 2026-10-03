package arm64

// Generator for sturb — one constructor (Sturb) over the shared
// load/store core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// SturbParams — parameters of the sturb form.
type SturbParams struct {
	LsParams
}

func NewSturbParams(p LsParams) SturbParams {
	return SturbParams{LsParams: p}
}

func (p SturbParams) Instr() arm64.Instr {
	in, err := arm64.New().Sturb(p.Rt, p.Rn, p.Off)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p SturbParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Sturb — an arbitrary sturb.
func Sturb(rnd *rand.Rand) ohsnap.Arbitrary[SturbParams] {
	rtW := false
	base := newLsGen(rnd, &rtW, UnscaledOff(rnd))
	return sturbArb{base: base}
}

type sturbArb struct {
	base lsGen
}

func (a sturbArb) Generate() iter.Seq[SturbParams] {
	return arbStream(func() SturbParams {
		return NewSturbParams(ohsnap.First(a.base.Generate()))
	})
}

func (a sturbArb) Shrink(p SturbParams) iter.Seq[SturbParams] {
	shrinks := slices.Collect(a.base.Shrink(p.LsParams))
	out := make([]SturbParams, 0, len(shrinks))
	for _, s := range shrinks {
		out = append(out, NewSturbParams(s))
	}

	return slices.Values(out)
}
