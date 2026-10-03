package arm64

// Generator for sdiv — one generator, one type, one constructor (Sdiv).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// SdivParams — parameters of sdiv rd, rn, rm.
type SdivParams struct {
	Rd, Rn, Rm arm64.Reg
}

func NewSdivParams(rd arm64.Reg, rn arm64.Reg, rm arm64.Reg) SdivParams {
	return SdivParams{
		Rd: rd,
		Rn: rn,
		Rm: rm,
	}
}

func (p SdivParams) Instr() arm64.Instr {
	in, err := arm64.New().Sdiv(p.Rd, p.Rn, p.Rm)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

func (p SdivParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// sdivGen — generator for sdiv: registers of the same width, 31st is zr.
type sdivGen struct {
	rnd *rand.Rand
}

// Sdiv — an arbitrary sdiv.
func Sdiv(rnd *rand.Rand) ohsnap.Arbitrary[SdivParams] {
	return newSdivGen(rnd)
}

func newSdivGen(rnd *rand.Rand) sdivGen {
	return sdivGen{rnd: rnd}
}

func (g sdivGen) Generate() iter.Seq[SdivParams] {
	return arb.Stream(func() SdivParams {
		is64 := g.rnd.IntN(2) == 1
		return NewSdivParams(
			genReg(g.rnd, is64, false, true),
			genReg(g.rnd, is64, false, true),
			genReg(g.rnd, is64, false, true),
		)
	})
}

func (g sdivGen) Shrink(p SdivParams) iter.Seq[SdivParams] {
	regs := regShrunk(p.Rd)
	out := make([]SdivParams, 0, len(regs))
	for _, r := range regs {
		out = append(out, NewSdivParams(r, p.Rn, p.Rm))
	}

	for _, r := range regShrunk(p.Rn) {
		out = append(out, NewSdivParams(p.Rd, r, p.Rm))
	}

	for _, r := range regShrunk(p.Rm) {
		out = append(out, NewSdivParams(p.Rd, p.Rn, r))
	}

	return slices.Values(out)
}
