package arm64

// Generator for adds (shifted register) — a thin signed variant over the
// AddShift family: one constructor (AddsShift), the generation and the
// shrink of the base add (shifted register).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// AddsShiftParams — parameters of adds rd, rn, rm[, shift #imm6].
type AddsShiftParams struct {
	Rd, Rn, Rm arm64.Reg
	Imm        arm64.Imm6
	Sh         arm64.Shift
}

func NewAddsShiftParams(rd arm64.Reg, rn arm64.Reg, rm arm64.Reg, imm arm64.Imm6, sh arm64.Shift) AddsShiftParams {
	return AddsShiftParams{
		Rd:  rd,
		Rn:  rn,
		Rm:  rm,
		Imm: imm,
		Sh:  sh,
	}
}

func (p AddsShiftParams) Instr() arm64.Instr {
	in, err := arm64.New().AddsShift(p.Rd, p.Rn, p.Rm, p.Imm, p.Sh)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p AddsShiftParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// AddsShift — an arbitrary adds (shifted register).
func AddsShift(rnd *rand.Rand) ohsnap.Arbitrary[AddsShiftParams] {
	base := AddShift(rnd)
	return addsShiftArb{base: base}
}

type addsShiftArb struct {
	base ohsnap.Arbitrary[AddShiftParams]
}

func (a addsShiftArb) Generate() iter.Seq[AddsShiftParams] {
	return arbStream(func() AddsShiftParams {
		return AddsShiftParams(ohsnap.First(a.base.Generate()))
	})
}

func (a addsShiftArb) Shrink(p AddsShiftParams) iter.Seq[AddsShiftParams] {
	var out []AddsShiftParams
	for _, s := range slices.Collect(a.base.Shrink(AddShiftParams(p))) {
		out = append(out, AddsShiftParams(s))
	}

	return slices.Values(out)
}
