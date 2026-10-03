package arm64

// Generator for sqrshl — one constructor (Sqrshl) over the shared
// three-register vector core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// SqrshlParams — parameters of sqrshl.
type SqrshlParams struct {
	V3Params
}

func NewSqrshlParams(p V3Params) SqrshlParams {
	return SqrshlParams{V3Params: p}
}

func (p SqrshlParams) Instr() arm64.Instr {
	in, err := arm64.New().Sqrshl(p.Rd, p.Rn, p.Rm, p.Arr)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p SqrshlParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Sqrshl — an arbitrary sqrshl.
func Sqrshl(rnd *rand.Rand) ohsnap.Arbitrary[SqrshlParams] {
	base := newV3Gen(rnd, arrFull())
	return sqrshlArb{base: base}
}

type sqrshlArb struct {
	base v3Gen
}

func (a sqrshlArb) Generate() iter.Seq[SqrshlParams] {
	return arbStream(func() SqrshlParams {
		return NewSqrshlParams(ohsnap.First(a.base.Generate()))
	})
}

func (a sqrshlArb) Shrink(p SqrshlParams) iter.Seq[SqrshlParams] {
	shrinks := slices.Collect(a.base.Shrink(p.V3Params))
	out := make([]SqrshlParams, 0, len(shrinks))
	for _, s := range shrinks {
		out = append(out, NewSqrshlParams(s))
	}

	return slices.Values(out)
}
