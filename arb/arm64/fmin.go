package arm64

// Generator for fmin — one constructor (Fmin) over the shared
// three-register FP core (one s or d kind).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// FminParams — parameters of fmin.
type FminParams struct {
	F3Params
}

func NewFminParams(p F3Params) FminParams {
	return FminParams{F3Params: p}
}

func (p FminParams) Instr() arm64.Instr {
	in, err := arm64.New().Fmin(p.Rd, p.Rn, p.Rm)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p FminParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Fmin — an arbitrary fmin.
func Fmin(rnd *rand.Rand) ohsnap.Arbitrary[FminParams] {
	base := f3(rnd)
	return fminArb{base: base}
}

type fminArb struct {
	base f3Gen
}

func (a fminArb) Generate() iter.Seq[FminParams] {
	return arbStream(func() FminParams {
		return NewFminParams(ohsnap.First(a.base.Generate()))
	})
}

func (a fminArb) Shrink(p FminParams) iter.Seq[FminParams] {
	shrinks := slices.Collect(a.base.Shrink(p.F3Params))
	out := make([]FminParams, 0, len(shrinks))
	for _, s := range shrinks {
		out = append(out, NewFminParams(s))
	}

	return slices.Values(out)
}
