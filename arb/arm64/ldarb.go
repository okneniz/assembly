package arm64

// Generator for ldarb — one constructor (Ldarb) over the shared
// atomic core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// LdarbParams — parameters of the ldarb form.
type LdarbParams struct {
	LsParams
}

func NewLdarbParams(p LsParams) LdarbParams {
	return LdarbParams{LsParams: p}
}

func (p LdarbParams) Instr() arm64.Instr {
	in, err := arm64.New().Ldarb(p.Rt, p.Rn)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p LdarbParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Ldarb — an arbitrary ldarb.
func Ldarb(rnd *rand.Rand) ohsnap.Arbitrary[LdarbParams] {
	base := atomW(rnd)
	return ldarbArb{base: base}
}

type ldarbArb struct {
	base atomGen
}

func (a ldarbArb) Generate() iter.Seq[LdarbParams] {
	return arbStream(func() LdarbParams {
		return NewLdarbParams(ohsnap.First(a.base.Generate()))
	})
}

func (a ldarbArb) Shrink(p LdarbParams) iter.Seq[LdarbParams] {
	shrinks := slices.Collect(a.base.Shrink(p.LsParams))
	out := make([]LdarbParams, 0, len(shrinks))
	for _, s := range shrinks {
		out = append(out, NewLdarbParams(s))
	}

	return slices.Values(out)
}
