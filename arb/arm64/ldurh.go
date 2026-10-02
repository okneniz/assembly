package arm64

// Generator for ldurh — one constructor (Ldurh) over the shared
// load/store core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// LdurhParams — parameters of the ldurh form.
type LdurhParams struct {
	LsParams
}

func NewLdurhParams(p LsParams) LdurhParams {
	return LdurhParams{LsParams: p}
}

func (p LdurhParams) Instr() arm64.Instr {
	in, err := arm64.New().Ldurh(p.Rt, p.Rn, p.Off)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p LdurhParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Ldurh — an arbitrary ldurh.
func Ldurh(rnd *rand.Rand) ohsnap.Arbitrary[LdurhParams] {
	rtW := false
	base := newLsGen(rnd, &rtW, UnscaledOff(rnd))
	return ldurhArb{base: base}
}

type ldurhArb struct {
	base lsGen
}

func (a ldurhArb) Generate() iter.Seq[LdurhParams] {
	return arbStream(func() LdurhParams {
		return NewLdurhParams(ohsnap.First(a.base.Generate()))
	})
}

func (a ldurhArb) Shrink(p LdurhParams) iter.Seq[LdurhParams] {
	var out []LdurhParams
	for _, s := range slices.Collect(a.base.Shrink(p.LsParams)) {
		out = append(out, NewLdurhParams(s))
	}

	return slices.Values(out)
}
