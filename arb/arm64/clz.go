package arm64

// Generator for clz — one generator, one type, one constructor (Clz).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// ClzParams — parameters of clz rd, rn.
type ClzParams struct {
	Rd, Rn arm64.Reg
}

func NewClzParams(rd arm64.Reg, rn arm64.Reg) ClzParams {
	return ClzParams{
		Rd: rd,
		Rn: rn,
	}
}

func (p ClzParams) Instr() arm64.Instr {
	in, err := arm64.New().Clz(p.Rd, p.Rn)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

func (p ClzParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// clzGen — generator for clz: registers of the same width, 31st is zr.
type clzGen struct {
	rnd *rand.Rand
}

// Clz — an arbitrary clz.
func Clz(rnd *rand.Rand) ohsnap.Arbitrary[ClzParams] {
	return newClzGen(rnd)
}

func newClzGen(rnd *rand.Rand) clzGen {
	return clzGen{rnd: rnd}
}

func (g clzGen) Generate() iter.Seq[ClzParams] {
	return arb.Stream(func() ClzParams {
		is64 := g.rnd.IntN(2) == 1
		return NewClzParams(genReg(g.rnd, is64, false, true), genReg(g.rnd, is64, false, true))
	})
}

func (g clzGen) Shrink(p ClzParams) iter.Seq[ClzParams] {
	regs := regShrunk(p.Rd)
	out := make([]ClzParams, 0, len(regs))
	for _, r := range regs {
		out = append(out, NewClzParams(r, p.Rn))
	}

	for _, r := range regShrunk(p.Rn) {
		out = append(out, NewClzParams(p.Rd, r))
	}

	return slices.Values(out)
}
