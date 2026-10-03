package arm64

// Generator for the ands shifted-register form — one constructor
// (AndsShift) over the shared shifted core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// AndsShiftParams — parameters of ands rd, rn, rm, shift #imm6.
type AndsShiftParams struct {
	ShiftedParams
}

func NewAndsShiftParams(p ShiftedParams) AndsShiftParams {
	return AndsShiftParams{ShiftedParams: p}
}

func (p AndsShiftParams) Instr() arm64.Instr {
	in, err := arm64.New().AndsShift(p.Rd, p.Rn, p.Rm, p.Imm6, p.Sh)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p AndsShiftParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// AndsShift — an arbitrary ands (shifted register).
func AndsShift(rnd *rand.Rand) ohsnap.Arbitrary[AndsShiftParams] {
	base := shifted(rnd)
	return andsShiftArb{base: base}
}

type andsShiftArb struct {
	base shiftedGen
}

func (a andsShiftArb) Generate() iter.Seq[AndsShiftParams] {
	return arbStream(func() AndsShiftParams {
		return NewAndsShiftParams(ohsnap.First(a.base.Generate()))
	})
}

func (a andsShiftArb) Shrink(p AndsShiftParams) iter.Seq[AndsShiftParams] {
	shrinks := slices.Collect(a.base.Shrink(p.ShiftedParams))
	out := make([]AndsShiftParams, 0, len(shrinks))
	for _, s := range shrinks {
		out = append(out, NewAndsShiftParams(s))
	}

	return slices.Values(out)
}
