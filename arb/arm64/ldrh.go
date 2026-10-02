package arm64

// Generator for ldrh — one constructor (Ldrh) over the shared
// load/store core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// LdrhParams — parameters of the ldrh form.
type LdrhParams struct {
	LsParams
}

func NewLdrhParams(p LsParams) LdrhParams {
	return LdrhParams{LsParams: p}
}

func (p LdrhParams) Instr() arm64.Instr {
	in, err := arm64.New().Ldrh(p.Rt, p.Rn, p.Off)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p LdrhParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Ldrh — an arbitrary ldrh.
func Ldrh(rnd *rand.Rand) ohsnap.Arbitrary[LdrhParams] {
	rtW := false
	base := newLsGen(rnd, &rtW, ScaledOff(rnd, 1))
	return ldrhArb{base: base}
}

type ldrhArb struct {
	base lsGen
}

func (a ldrhArb) Generate() iter.Seq[LdrhParams] {
	return arbStream(func() LdrhParams {
		return NewLdrhParams(ohsnap.First(a.base.Generate()))
	})
}

func (a ldrhArb) Shrink(p LdrhParams) iter.Seq[LdrhParams] {
	var out []LdrhParams
	for _, s := range slices.Collect(a.base.Shrink(p.LsParams)) {
		out = append(out, NewLdrhParams(s))
	}

	return slices.Values(out)
}
