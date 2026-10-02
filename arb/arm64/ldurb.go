package arm64

// Generator for ldurb — one constructor (Ldurb) over the shared
// load/store core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// LdurbParams — parameters of the ldurb form.
type LdurbParams struct {
	LsParams
}

func NewLdurbParams(p LsParams) LdurbParams {
	return LdurbParams{LsParams: p}
}

func (p LdurbParams) Instr() arm64.Instr {
	in, err := arm64.New().Ldurb(p.Rt, p.Rn, p.Off)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p LdurbParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Ldurb — an arbitrary ldurb.
func Ldurb(rnd *rand.Rand) ohsnap.Arbitrary[LdurbParams] {
	rtW := false
	base := newLsGen(rnd, &rtW, UnscaledOff(rnd))
	return ldurbArb{base: base}
}

type ldurbArb struct {
	base lsGen
}

func (a ldurbArb) Generate() iter.Seq[LdurbParams] {
	return arbStream(func() LdurbParams {
		return NewLdurbParams(ohsnap.First(a.base.Generate()))
	})
}

func (a ldurbArb) Shrink(p LdurbParams) iter.Seq[LdurbParams] {
	var out []LdurbParams
	for _, s := range slices.Collect(a.base.Shrink(p.LsParams)) {
		out = append(out, NewLdurbParams(s))
	}

	return slices.Values(out)
}
