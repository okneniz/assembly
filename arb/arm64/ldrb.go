package arm64

// Generator for ldrb — one constructor (Ldrb) over the shared
// load/store core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// LdrbParams — parameters of the ldrb form.
type LdrbParams struct {
	LsParams
}

func NewLdrbParams(p LsParams) LdrbParams {
	return LdrbParams{LsParams: p}
}

func (p LdrbParams) Instr() arm64.Instr {
	in, err := arm64.New().Ldrb(p.Rt, p.Rn, p.Off)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p LdrbParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Ldrb — an arbitrary ldrb.
func Ldrb(rnd *rand.Rand) ohsnap.Arbitrary[LdrbParams] {
	rtW := false
	base := newLsGen(rnd, &rtW, ScaledOff(rnd, 0))
	return ldrbArb{base: base}
}

type ldrbArb struct {
	base lsGen
}

func (a ldrbArb) Generate() iter.Seq[LdrbParams] {
	return arbStream(func() LdrbParams {
		return NewLdrbParams(ohsnap.First(a.base.Generate()))
	})
}

func (a ldrbArb) Shrink(p LdrbParams) iter.Seq[LdrbParams] {
	shrinks := slices.Collect(a.base.Shrink(p.LsParams))
	out := make([]LdrbParams, 0, len(shrinks))
	for _, s := range shrinks {
		out = append(out, NewLdrbParams(s))
	}

	return slices.Values(out)
}
