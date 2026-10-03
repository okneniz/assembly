package arm64

// Generator for the bic shifted-register form — one constructor
// (BicShift) over the shared shifted core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// BicShiftParams — parameters of bic rd, rn, rm, shift #imm6.
type BicShiftParams struct {
	ShiftedParams
}

func NewBicShiftParams(p ShiftedParams) BicShiftParams {
	return BicShiftParams{ShiftedParams: p}
}

func (p BicShiftParams) Instr() arm64.Instr {
	in, err := arm64.New().BicShift(p.Rd, p.Rn, p.Rm, p.Imm6, p.Sh)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p BicShiftParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// BicShift — an arbitrary bic (shifted register).
func BicShift(rnd *rand.Rand) ohsnap.Arbitrary[BicShiftParams] {
	base := shifted(rnd)
	return bicShiftArb{base: base}
}

type bicShiftArb struct {
	base shiftedGen
}

func (a bicShiftArb) Generate() iter.Seq[BicShiftParams] {
	return arbStream(func() BicShiftParams {
		return NewBicShiftParams(ohsnap.First(a.base.Generate()))
	})
}

func (a bicShiftArb) Shrink(p BicShiftParams) iter.Seq[BicShiftParams] {
	shrinks := slices.Collect(a.base.Shrink(p.ShiftedParams))
	out := make([]BicShiftParams, 0, len(shrinks))
	for _, s := range shrinks {
		out = append(out, NewBicShiftParams(s))
	}

	return slices.Values(out)
}
