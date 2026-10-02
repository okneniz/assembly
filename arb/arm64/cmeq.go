package arm64

// Generator for cmeq — one constructor (Cmeq) over the shared
// three-register vector core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// CmeqParams — parameters of cmeq.
type CmeqParams struct {
	V3Params
}

func NewCmeqParams(p V3Params) CmeqParams {
	return CmeqParams{V3Params: p}
}

func (p CmeqParams) Instr() arm64.Instr {
	in, err := arm64.New().Cmeq(p.Rd, p.Rn, p.Rm, p.Arr)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p CmeqParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Cmeq — an arbitrary cmeq.
func Cmeq(rnd *rand.Rand) ohsnap.Arbitrary[CmeqParams] {
	base := newV3Gen(rnd, arrFull())
	return cmeqArb{base: base}
}

type cmeqArb struct {
	base v3Gen
}

func (a cmeqArb) Generate() iter.Seq[CmeqParams] {
	return arbStream(func() CmeqParams {
		return NewCmeqParams(ohsnap.First(a.base.Generate()))
	})
}

func (a cmeqArb) Shrink(p CmeqParams) iter.Seq[CmeqParams] {
	var out []CmeqParams
	for _, s := range slices.Collect(a.base.Shrink(p.V3Params)) {
		out = append(out, NewCmeqParams(s))
	}

	return slices.Values(out)
}
