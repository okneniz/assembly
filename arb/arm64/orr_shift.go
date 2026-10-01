package arm64

// Generator for the orr shifted-register form — one constructor
// (OrrShift) over the shared shifted core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// OrrShiftParams — parameters of orr rd, rn, rm, shift #imm6.
type OrrShiftParams struct {
	ShiftedParams
}

func NewOrrShiftParams(p ShiftedParams) OrrShiftParams {
	return OrrShiftParams{ShiftedParams: p}
}

func (p OrrShiftParams) Instr() arm64.Instr {
	in, err := arm64.New().OrrShift(p.Rd, p.Rn, p.Rm, p.Imm6, p.Sh)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p OrrShiftParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// OrrShift — an arbitrary orr (shifted register).
func OrrShift(rnd *rand.Rand) ohsnap.Arbitrary[OrrShiftParams] {
	base := shifted(rnd)
	return orrShiftArb{base: base}
}

type orrShiftArb struct {
	base shiftedGen
}

func (a orrShiftArb) Generate() iter.Seq[OrrShiftParams] {
	return arbStream(func() OrrShiftParams {
		return NewOrrShiftParams(ohsnap.First(a.base.Generate()))
	})
}

func (a orrShiftArb) Shrink(p OrrShiftParams) iter.Seq[OrrShiftParams] {
	var out []OrrShiftParams
	for _, s := range slices.Collect(a.base.Shrink(p.ShiftedParams)) {
		out = append(out, NewOrrShiftParams(s))
	}

	return slices.Values(out)
}
