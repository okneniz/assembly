package arm64

// Generator for ldaxr — one constructor (Ldaxr) over the shared
// atomic core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// LdaxrParams — parameters of the ldaxr form.
type LdaxrParams struct {
	LsParams
}

func NewLdaxrParams(p LsParams) LdaxrParams {
	return LdaxrParams{LsParams: p}
}

func (p LdaxrParams) Instr() arm64.Instr {
	in, err := arm64.New().Ldaxr(p.Rt, p.Rn)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p LdaxrParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Ldaxr — an arbitrary ldaxr.
func Ldaxr(rnd *rand.Rand) ohsnap.Arbitrary[LdaxrParams] {
	base := atom(rnd)
	return ldaxrArb{base: base}
}

type ldaxrArb struct {
	base atomGen
}

func (a ldaxrArb) Generate() iter.Seq[LdaxrParams] {
	return arbStream(func() LdaxrParams {
		return NewLdaxrParams(ohsnap.First(a.base.Generate()))
	})
}

func (a ldaxrArb) Shrink(p LdaxrParams) iter.Seq[LdaxrParams] {
	shrinks := slices.Collect(a.base.Shrink(p.LsParams))
	out := make([]LdaxrParams, 0, len(shrinks))
	for _, s := range shrinks {
		out = append(out, NewLdaxrParams(s))
	}

	return slices.Values(out)
}
