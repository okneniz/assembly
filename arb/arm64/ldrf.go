package arm64

// Generator for ldr (FP register form) — one constructor (LdrF) over the
// shared f core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// LdrFParams — parameters of ldr st|dt, [xn, #off].
type LdrFParams struct {
	FParams
}

func NewLdrFParams(p FParams) LdrFParams {
	return LdrFParams{FParams: p}
}

func (p LdrFParams) Instr() arm64.Instr {
	in, err := arm64.New().LdrF(p.Rt, p.Rn, p.Off)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p LdrFParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// LdrF — an arbitrary ldr (FP register).
func LdrF(rnd *rand.Rand) ohsnap.Arbitrary[LdrFParams] {
	base := f(rnd)
	return ldrFArb{base: base}
}

type ldrFArb struct {
	base fGen
}

func (a ldrFArb) Generate() iter.Seq[LdrFParams] {
	return arbStream(func() LdrFParams {
		return NewLdrFParams(ohsnap.First(a.base.Generate()))
	})
}

func (a ldrFArb) Shrink(p LdrFParams) iter.Seq[LdrFParams] {
	shrinks := slices.Collect(a.base.Shrink(p.FParams))
	out := make([]LdrFParams, 0, len(shrinks))
	for _, s := range shrinks {
		out = append(out, NewLdrFParams(s))
	}

	return slices.Values(out)
}
