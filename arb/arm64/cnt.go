package arm64

// Generator for cnt — one constructor (Cnt) over the shared
// two-register vector core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// CntParams — parameters of cnt.
type CntParams struct {
	V2Params
}

func NewCntParams(p V2Params) CntParams {
	return CntParams{V2Params: p}
}

func (p CntParams) Instr() arm64.Instr {
	in, err := arm64.New().Cnt(p.Rd, p.Rn, p.Arr)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p CntParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Cnt — an arbitrary cnt.
func Cnt(rnd *rand.Rand) ohsnap.Arbitrary[CntParams] {
	base := newV2Gen(rnd, arrLogical())
	return cntArb{base: base}
}

type cntArb struct {
	base v2Gen
}

func (a cntArb) Generate() iter.Seq[CntParams] {
	return arbStream(func() CntParams {
		return NewCntParams(ohsnap.First(a.base.Generate()))
	})
}

func (a cntArb) Shrink(p CntParams) iter.Seq[CntParams] {
	var out []CntParams
	for _, s := range slices.Collect(a.base.Shrink(p.V2Params)) {
		out = append(out, NewCntParams(s))
	}

	return slices.Values(out)
}
