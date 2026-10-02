package arm64

// Generator for str (FP register form) — one constructor (StrF) over the
// shared f core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// StrFParams — parameters of str st|dt, [xn, #off].
type StrFParams struct {
	FParams
}

func NewStrFParams(p FParams) StrFParams {
	return StrFParams{FParams: p}
}

func (p StrFParams) Instr() arm64.Instr {
	in, err := arm64.New().StrF(p.Rt, p.Rn, p.Off)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p StrFParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// StrF — an arbitrary str (FP register).
func StrF(rnd *rand.Rand) ohsnap.Arbitrary[StrFParams] {
	base := f(rnd)
	return strFArb{base: base}
}

type strFArb struct {
	base fGen
}

func (a strFArb) Generate() iter.Seq[StrFParams] {
	return arbStream(func() StrFParams {
		return NewStrFParams(ohsnap.First(a.base.Generate()))
	})
}

func (a strFArb) Shrink(p StrFParams) iter.Seq[StrFParams] {
	var out []StrFParams
	for _, s := range slices.Collect(a.base.Shrink(p.FParams)) {
		out = append(out, NewStrFParams(s))
	}

	return slices.Values(out)
}
