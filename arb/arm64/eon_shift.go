package arm64

// Generator for the eon shifted-register form — one constructor
// (EonShift) over the shared shifted core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// EonShiftParams — parameters of eon rd, rn, rm, shift #imm6.
type EonShiftParams struct {
	ShiftedParams
}

func NewEonShiftParams(p ShiftedParams) EonShiftParams {
	return EonShiftParams{ShiftedParams: p}
}

func (p EonShiftParams) Instr() arm64.Instr {
	in, err := arm64.New().EonShift(p.Rd, p.Rn, p.Rm, p.Imm6, p.Sh)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p EonShiftParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// EonShift — an arbitrary eon (shifted register).
func EonShift(rnd *rand.Rand) ohsnap.Arbitrary[EonShiftParams] {
	base := shifted(rnd)
	return eonShiftArb{base: base}
}

type eonShiftArb struct {
	base shiftedGen
}

func (a eonShiftArb) Generate() iter.Seq[EonShiftParams] {
	return arbStream(func() EonShiftParams {
		return NewEonShiftParams(ohsnap.First(a.base.Generate()))
	})
}

func (a eonShiftArb) Shrink(p EonShiftParams) iter.Seq[EonShiftParams] {
	shrinks := slices.Collect(a.base.Shrink(p.ShiftedParams))
	out := make([]EonShiftParams, 0, len(shrinks))
	for _, s := range shrinks {
		out = append(out, NewEonShiftParams(s))
	}

	return slices.Values(out)
}
