package arm64

// Generator for eor — one constructor (Eor) over the shared
// three-register vector core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// EorVParams — parameters of eor.
type EorVParams struct {
	V3Params
}

func NewEorVParams(p V3Params) EorVParams {
	return EorVParams{V3Params: p}
}

func (p EorVParams) Instr() arm64.Instr {
	in, err := arm64.New().Eor(p.Rd, p.Rn, p.Rm, p.Arr)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p EorVParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Eor — an arbitrary eor.
func Eor(rnd *rand.Rand) ohsnap.Arbitrary[EorVParams] {
	base := newV3Gen(rnd, arrLogical())
	return eorArb{base: base}
}

type eorArb struct {
	base v3Gen
}

func (a eorArb) Generate() iter.Seq[EorVParams] {
	return arbStream(func() EorVParams {
		return NewEorVParams(ohsnap.First(a.base.Generate()))
	})
}

func (a eorArb) Shrink(p EorVParams) iter.Seq[EorVParams] {
	var out []EorVParams
	for _, s := range slices.Collect(a.base.Shrink(p.V3Params)) {
		out = append(out, NewEorVParams(s))
	}

	return slices.Values(out)
}
