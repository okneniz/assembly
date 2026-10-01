package arm64

// Generator for adds (immediate) — a thin signed variant over the AddImm
// family: one constructor (AddsImm), the generation and the shrink of the
// base add (immediate).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// AddsImmParams — parameters of adds rd, rn, #imm12[, lsl #12].
type AddsImmParams struct {
	Rd, Rn arm64.Reg
	Imm    arm64.Imm12
	Sh     arm64.Sh12
}

func NewAddsImmParams(rd arm64.Reg, rn arm64.Reg, imm arm64.Imm12, sh arm64.Sh12) AddsImmParams {
	return AddsImmParams{
		Rd:  rd,
		Rn:  rn,
		Imm: imm,
		Sh:  sh,
	}
}

func (p AddsImmParams) Instr() arm64.Instr {
	in, err := arm64.New().AddsImm(p.Rd, p.Rn, p.Imm, p.Sh)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p AddsImmParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// AddsImm — an arbitrary adds (immediate).
func AddsImm(rnd *rand.Rand) ohsnap.Arbitrary[AddsImmParams] {
	base := AddImm(rnd)
	return addsImmArb{base: base}
}

type addsImmArb struct {
	base ohsnap.Arbitrary[AddImmParams]
}

func (a addsImmArb) Generate() iter.Seq[AddsImmParams] {
	return arbStream(func() AddsImmParams {
		return AddsImmParams(ohsnap.First(a.base.Generate()))
	})
}

func (a addsImmArb) Shrink(p AddsImmParams) iter.Seq[AddsImmParams] {
	var out []AddsImmParams
	for _, s := range slices.Collect(a.base.Shrink(AddImmParams(p))) {
		out = append(out, AddsImmParams(s))
	}

	return slices.Values(out)
}
