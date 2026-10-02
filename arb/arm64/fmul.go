package arm64

// Generator for fmul — one constructor (Fmul) over the shared
// three-register FP core (one s or d kind).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// FmulParams — parameters of fmul.
type FmulParams struct {
	F3Params
}

func NewFmulParams(p F3Params) FmulParams {
	return FmulParams{F3Params: p}
}

func (p FmulParams) Instr() arm64.Instr {
	in, err := arm64.New().Fmul(p.Rd, p.Rn, p.Rm)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p FmulParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Fmul — an arbitrary fmul.
func Fmul(rnd *rand.Rand) ohsnap.Arbitrary[FmulParams] {
	base := f3(rnd)
	return fmulArb{base: base}
}

type fmulArb struct {
	base f3Gen
}

func (a fmulArb) Generate() iter.Seq[FmulParams] {
	return arbStream(func() FmulParams {
		return NewFmulParams(ohsnap.First(a.base.Generate()))
	})
}

func (a fmulArb) Shrink(p FmulParams) iter.Seq[FmulParams] {
	var out []FmulParams
	for _, s := range slices.Collect(a.base.Shrink(p.F3Params)) {
		out = append(out, NewFmulParams(s))
	}

	return slices.Values(out)
}
