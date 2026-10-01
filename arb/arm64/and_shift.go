package arm64

// Generator for the and shifted-register form — one constructor
// (AndShift) over the shared shifted core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// AndShiftParams — parameters of and rd, rn, rm, shift #imm6.
type AndShiftParams struct {
	ShiftedParams
}

func NewAndShiftParams(p ShiftedParams) AndShiftParams {
	return AndShiftParams{ShiftedParams: p}
}

func (p AndShiftParams) Instr() arm64.Instr {
	in, err := arm64.New().AndShift(p.Rd, p.Rn, p.Rm, p.Imm6, p.Sh)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p AndShiftParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// AndShift — an arbitrary and (shifted register).
func AndShift(rnd *rand.Rand) ohsnap.Arbitrary[AndShiftParams] {
	base := shifted(rnd)
	return andShiftArb{base: base}
}

type andShiftArb struct {
	base shiftedGen
}

func (a andShiftArb) Generate() iter.Seq[AndShiftParams] {
	return arbStream(func() AndShiftParams {
		return NewAndShiftParams(ohsnap.First(a.base.Generate()))
	})
}

func (a andShiftArb) Shrink(p AndShiftParams) iter.Seq[AndShiftParams] {
	var out []AndShiftParams
	for _, s := range slices.Collect(a.base.Shrink(p.ShiftedParams)) {
		out = append(out, NewAndShiftParams(s))
	}

	return slices.Values(out)
}
