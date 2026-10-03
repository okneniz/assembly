package arm64

// Generator for the orn shifted-register form — one constructor
// (OrnShift) over the shared shifted core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// OrnShiftParams — parameters of orn rd, rn, rm, shift #imm6.
type OrnShiftParams struct {
	ShiftedParams
}

func NewOrnShiftParams(p ShiftedParams) OrnShiftParams {
	return OrnShiftParams{ShiftedParams: p}
}

func (p OrnShiftParams) Instr() arm64.Instr {
	in, err := arm64.New().OrnShift(p.Rd, p.Rn, p.Rm, p.Imm6, p.Sh)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p OrnShiftParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// OrnShift — an arbitrary orn (shifted register).
func OrnShift(rnd *rand.Rand) ohsnap.Arbitrary[OrnShiftParams] {
	base := shifted(rnd)
	return ornShiftArb{base: base}
}

type ornShiftArb struct {
	base shiftedGen
}

func (a ornShiftArb) Generate() iter.Seq[OrnShiftParams] {
	return arbStream(func() OrnShiftParams {
		return NewOrnShiftParams(ohsnap.First(a.base.Generate()))
	})
}

func (a ornShiftArb) Shrink(p OrnShiftParams) iter.Seq[OrnShiftParams] {
	shrinks := slices.Collect(a.base.Shrink(p.ShiftedParams))
	out := make([]OrnShiftParams, 0, len(shrinks))
	for _, s := range shrinks {
		out = append(out, NewOrnShiftParams(s))
	}

	return slices.Values(out)
}
