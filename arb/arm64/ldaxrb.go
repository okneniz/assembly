package arm64

// Generator for ldaxrb — one constructor (Ldaxrb) over the shared
// atomic core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// LdaxrbParams — parameters of the ldaxrb form.
type LdaxrbParams struct {
	LsParams
}

func NewLdaxrbParams(p LsParams) LdaxrbParams {
	return LdaxrbParams{LsParams: p}
}

func (p LdaxrbParams) Instr() arm64.Instr {
	in, err := arm64.New().Ldaxrb(p.Rt, p.Rn)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p LdaxrbParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Ldaxrb — an arbitrary ldaxrb.
func Ldaxrb(rnd *rand.Rand) ohsnap.Arbitrary[LdaxrbParams] {
	base := atomW(rnd)
	return ldaxrbArb{base: base}
}

type ldaxrbArb struct {
	base atomGen
}

func (a ldaxrbArb) Generate() iter.Seq[LdaxrbParams] {
	return arbStream(func() LdaxrbParams {
		return NewLdaxrbParams(ohsnap.First(a.base.Generate()))
	})
}

func (a ldaxrbArb) Shrink(p LdaxrbParams) iter.Seq[LdaxrbParams] {
	var out []LdaxrbParams
	for _, s := range slices.Collect(a.base.Shrink(p.LsParams)) {
		out = append(out, NewLdaxrbParams(s))
	}

	return slices.Values(out)
}
