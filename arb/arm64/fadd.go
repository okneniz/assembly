package arm64

// Generator for fadd — one constructor (Fadd) over the shared
// three-register FP core (one s or d kind).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// FaddParams — parameters of fadd.
type FaddParams struct {
	F3Params
}

func NewFaddParams(p F3Params) FaddParams {
	return FaddParams{F3Params: p}
}

func (p FaddParams) Instr() arm64.Instr {
	in, err := arm64.New().Fadd(p.Rd, p.Rn, p.Rm)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p FaddParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Fadd — an arbitrary fadd.
func Fadd(rnd *rand.Rand) ohsnap.Arbitrary[FaddParams] {
	base := f3(rnd)
	return faddArb{base: base}
}

type faddArb struct {
	base f3Gen
}

func (a faddArb) Generate() iter.Seq[FaddParams] {
	return arbStream(func() FaddParams {
		return NewFaddParams(ohsnap.First(a.base.Generate()))
	})
}

func (a faddArb) Shrink(p FaddParams) iter.Seq[FaddParams] {
	var out []FaddParams
	for _, s := range slices.Collect(a.base.Shrink(p.F3Params)) {
		out = append(out, NewFaddParams(s))
	}

	return slices.Values(out)
}
