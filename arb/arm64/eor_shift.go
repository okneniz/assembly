package arm64

// Generator for the eor shifted-register form — one constructor
// (EorShift) over the shared shifted core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// EorShiftParams — parameters of eor rd, rn, rm, shift #imm6.
type EorShiftParams struct {
	ShiftedParams
}

func NewEorShiftParams(p ShiftedParams) EorShiftParams {
	return EorShiftParams{ShiftedParams: p}
}

func (p EorShiftParams) Instr() arm64.Instr {
	in, err := arm64.New().EorShift(p.Rd, p.Rn, p.Rm, p.Imm6, p.Sh)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p EorShiftParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// EorShift — an arbitrary eor (shifted register).
func EorShift(rnd *rand.Rand) ohsnap.Arbitrary[EorShiftParams] {
	base := shifted(rnd)
	return eorShiftArb{base: base}
}

type eorShiftArb struct {
	base shiftedGen
}

func (a eorShiftArb) Generate() iter.Seq[EorShiftParams] {
	return arbStream(func() EorShiftParams {
		return NewEorShiftParams(ohsnap.First(a.base.Generate()))
	})
}

func (a eorShiftArb) Shrink(p EorShiftParams) iter.Seq[EorShiftParams] {
	shrinks := slices.Collect(a.base.Shrink(p.ShiftedParams))
	out := make([]EorShiftParams, 0, len(shrinks))
	for _, s := range shrinks {
		out = append(out, NewEorShiftParams(s))
	}

	return slices.Values(out)
}
