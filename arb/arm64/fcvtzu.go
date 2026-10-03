package arm64

// Generator for fcvtzu — one constructor (Fcvtzu) over the shared
// conversion core (rd is the gpr, rn the FP register).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// FcvtzuParams — parameters of fcvtzu.
type FcvtzuParams struct {
	FGprParams
}

func NewFcvtzuParams(p FGprParams) FcvtzuParams {
	return FcvtzuParams{FGprParams: p}
}

func (p FcvtzuParams) Instr() arm64.Instr {
	in, err := arm64.New().Fcvtzu(p.R, p.F)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p FcvtzuParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Fcvtzu — an arbitrary fcvtzu.
func Fcvtzu(rnd *rand.Rand) ohsnap.Arbitrary[FcvtzuParams] {
	base := fgprAny(rnd)
	return fcvtzuArb{base: base}
}

type fcvtzuArb struct {
	base fgprGen
}

func (a fcvtzuArb) Generate() iter.Seq[FcvtzuParams] {
	return arbStream(func() FcvtzuParams {
		return NewFcvtzuParams(ohsnap.First(a.base.Generate()))
	})
}

func (a fcvtzuArb) Shrink(p FcvtzuParams) iter.Seq[FcvtzuParams] {
	shrinks := slices.Collect(a.base.Shrink(p.FGprParams))
	out := make([]FcvtzuParams, 0, len(shrinks))
	for _, s := range shrinks {
		out = append(out, NewFcvtzuParams(s))
	}

	return slices.Values(out)
}
