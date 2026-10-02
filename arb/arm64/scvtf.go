package arm64

// Generator for scvtf — one constructor (Scvtf) over the shared
// conversion core (rd is the FP register, rn the gpr).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// ScvtfParams — parameters of scvtf.
type ScvtfParams struct {
	FGprParams
}

func NewScvtfParams(p FGprParams) ScvtfParams {
	return ScvtfParams{FGprParams: p}
}

func (p ScvtfParams) Instr() arm64.Instr {
	in, err := arm64.New().Scvtf(p.F, p.R)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p ScvtfParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Scvtf — an arbitrary scvtf.
func Scvtf(rnd *rand.Rand) ohsnap.Arbitrary[ScvtfParams] {
	base := fgprAny(rnd)
	return scvtfArb{base: base}
}

type scvtfArb struct {
	base fgprGen
}

func (a scvtfArb) Generate() iter.Seq[ScvtfParams] {
	return arbStream(func() ScvtfParams {
		return NewScvtfParams(ohsnap.First(a.base.Generate()))
	})
}

func (a scvtfArb) Shrink(p ScvtfParams) iter.Seq[ScvtfParams] {
	var out []ScvtfParams
	for _, s := range slices.Collect(a.base.Shrink(p.FGprParams)) {
		out = append(out, NewScvtfParams(s))
	}

	return slices.Values(out)
}
