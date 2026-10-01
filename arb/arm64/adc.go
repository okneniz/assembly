package arm64

// Generator for adc — one generator, one type, one constructor (Adc).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// AdcParams — parameters of adc rd, rn, rm (the 64-bit form only).
type AdcParams struct {
	Rd, Rn, Rm arm64.Reg
}

func NewAdcParams(rd arm64.Reg, rn arm64.Reg, rm arm64.Reg) AdcParams {
	return AdcParams{
		Rd: rd,
		Rn: rn,
		Rm: rm,
	}
}

func (p AdcParams) Instr() arm64.Instr {
	in, err := arm64.New().Adc(p.Rd, p.Rn, p.Rm)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p AdcParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// adcGen — generator for adc: x-registers, occasionally xzr.
type adcGen struct {
	rnd *rand.Rand
}

func newAdcGen(rnd *rand.Rand) adcGen {
	return adcGen{rnd: rnd}
}

// Adc — an arbitrary adc.
func Adc(rnd *rand.Rand) ohsnap.Arbitrary[AdcParams] {
	return newAdcGen(rnd)
}

func (g adcGen) Generate() iter.Seq[AdcParams] {
	return arb.Stream(func() AdcParams {
		return NewAdcParams(
			genReg(g.rnd, true, false, true),
			genReg(g.rnd, true, false, true),
			genReg(g.rnd, true, false, true),
		)
	})
}

func (g adcGen) Shrink(p AdcParams) iter.Seq[AdcParams] {
	var out []AdcParams
	for _, r := range regShrunk(p.Rd) {
		out = append(out, NewAdcParams(r, p.Rn, p.Rm))
	}

	for _, r := range regShrunk(p.Rn) {
		out = append(out, NewAdcParams(p.Rd, r, p.Rm))
	}

	for _, r := range regShrunk(p.Rm) {
		out = append(out, NewAdcParams(p.Rd, p.Rn, r))
	}

	return slices.Values(out)
}
