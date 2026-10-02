package arm64

// Generator for uaddw — one constructor (Uaddw) over the shared
// three-register vector core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// UaddwParams — parameters of uaddw.
type UaddwParams struct {
	V3Params
}

func NewUaddwParams(p V3Params) UaddwParams {
	return UaddwParams{V3Params: p}
}

func (p UaddwParams) Instr() arm64.Instr {
	in, err := arm64.New().Uaddw(p.Rd, p.Rn, p.Rm, p.Arr)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p UaddwParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Uaddw — an arbitrary uaddw.
func Uaddw(rnd *rand.Rand) ohsnap.Arbitrary[UaddwParams] {
	base := newV3Gen(rnd, arrWiden())
	return uaddwArb{base: base}
}

type uaddwArb struct {
	base v3Gen
}

func (a uaddwArb) Generate() iter.Seq[UaddwParams] {
	return arbStream(func() UaddwParams {
		return NewUaddwParams(ohsnap.First(a.base.Generate()))
	})
}

func (a uaddwArb) Shrink(p UaddwParams) iter.Seq[UaddwParams] {
	var out []UaddwParams
	for _, s := range slices.Collect(a.base.Shrink(p.V3Params)) {
		out = append(out, NewUaddwParams(s))
	}

	return slices.Values(out)
}
