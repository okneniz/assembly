package arm64

// Generator for the bics shifted-register form — one constructor
// (BicsShift) over the shared shifted core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// BicsShiftParams — parameters of bics rd, rn, rm, shift #imm6.
type BicsShiftParams struct {
	ShiftedParams
}

func NewBicsShiftParams(p ShiftedParams) BicsShiftParams {
	return BicsShiftParams{ShiftedParams: p}
}

func (p BicsShiftParams) Instr() arm64.Instr {
	in, err := arm64.New().BicsShift(p.Rd, p.Rn, p.Rm, p.Imm6, p.Sh)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p BicsShiftParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// BicsShift — an arbitrary bics (shifted register).
func BicsShift(rnd *rand.Rand) ohsnap.Arbitrary[BicsShiftParams] {
	base := shifted(rnd)
	return bicsShiftArb{base: base}
}

type bicsShiftArb struct {
	base shiftedGen
}

func (a bicsShiftArb) Generate() iter.Seq[BicsShiftParams] {
	return arbStream(func() BicsShiftParams {
		return NewBicsShiftParams(ohsnap.First(a.base.Generate()))
	})
}

func (a bicsShiftArb) Shrink(p BicsShiftParams) iter.Seq[BicsShiftParams] {
	var out []BicsShiftParams
	for _, s := range slices.Collect(a.base.Shrink(p.ShiftedParams)) {
		out = append(out, NewBicsShiftParams(s))
	}

	return slices.Values(out)
}
