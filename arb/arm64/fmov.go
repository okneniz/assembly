package arm64

// Generator for fmov — one constructor (Fmov) over the shared
// two-register FP core (one s or d kind).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// FmovParams — parameters of fmov.
type FmovParams struct {
	F2Params
}

func NewFmovParams(p F2Params) FmovParams {
	return FmovParams{F2Params: p}
}

func (p FmovParams) Instr() arm64.Instr {
	in, err := arm64.New().Fmov(p.Rd, p.Rn)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p FmovParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Fmov — an arbitrary fmov.
func Fmov(rnd *rand.Rand) ohsnap.Arbitrary[FmovParams] {
	base := f2(rnd)
	return fmovArb{base: base}
}

type fmovArb struct {
	base f2Gen
}

func (a fmovArb) Generate() iter.Seq[FmovParams] {
	return arbStream(func() FmovParams {
		return NewFmovParams(ohsnap.First(a.base.Generate()))
	})
}

func (a fmovArb) Shrink(p FmovParams) iter.Seq[FmovParams] {
	var out []FmovParams
	for _, s := range slices.Collect(a.base.Shrink(p.F2Params)) {
		out = append(out, NewFmovParams(s))
	}

	return slices.Values(out)
}
