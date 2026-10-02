package arm64

// Generator for uaddlv — one constructor (Uaddlv) over the shared
// two-register vector core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// UaddlvParams — parameters of uaddlv.
type UaddlvParams struct {
	V2Params
}

func NewUaddlvParams(p V2Params) UaddlvParams {
	return UaddlvParams{V2Params: p}
}

func (p UaddlvParams) Instr() arm64.Instr {
	in, err := arm64.New().Uaddlv(p.Rd, p.Rn, p.Arr)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p UaddlvParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Uaddlv — an arbitrary uaddlv.
func Uaddlv(rnd *rand.Rand) ohsnap.Arbitrary[UaddlvParams] {
	base := newV2Gen(rnd, arrHalf())
	return uaddlvArb{base: base}
}

type uaddlvArb struct {
	base v2Gen
}

func (a uaddlvArb) Generate() iter.Seq[UaddlvParams] {
	return arbStream(func() UaddlvParams {
		return NewUaddlvParams(ohsnap.First(a.base.Generate()))
	})
}

func (a uaddlvArb) Shrink(p UaddlvParams) iter.Seq[UaddlvParams] {
	var out []UaddlvParams
	for _, s := range slices.Collect(a.base.Shrink(p.V2Params)) {
		out = append(out, NewUaddlvParams(s))
	}

	return slices.Values(out)
}
