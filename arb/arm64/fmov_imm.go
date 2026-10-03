package arm64

// Generator for fmov (immediate) — one constructor (FmovImm) over the
// imm8 space: every candidate expands through the arch VfpExpand tables,
// so the value is encodable by construction.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// FmovImmParams — parameters of the fmov immediate form.
type FmovImmParams struct {
	FpImmParams
}

func NewFmovImmParams(p FpImmParams) FmovImmParams {
	return FmovImmParams{FpImmParams: p}
}

func (p FmovImmParams) Instr() arm64.Instr {
	in, err := arm64.New().FmovImm(p.Rd, p.Value())
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p FmovImmParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// FmovImm — an arbitrary fmov (immediate).
func FmovImm(rnd *rand.Rand) ohsnap.Arbitrary[FmovImmParams] {
	base := newFpImmGen(rnd)
	return fmovImmArb{base: base}
}

type fmovImmArb struct {
	base fpImmGen
}

func (a fmovImmArb) Generate() iter.Seq[FmovImmParams] {
	return arbStream(func() FmovImmParams {
		return NewFmovImmParams(ohsnap.First(a.base.Generate()))
	})
}

func (a fmovImmArb) Shrink(p FmovImmParams) iter.Seq[FmovImmParams] {
	shrinks := slices.Collect(a.base.Shrink(p.FpImmParams))
	out := make([]FmovImmParams, 0, len(shrinks))
	for _, s := range shrinks {
		out = append(out, NewFmovImmParams(s))
	}

	return slices.Values(out)
}
