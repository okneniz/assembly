package arm64

// Generator for cls — one generator, one type, one constructor (Cls).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// ClsParams — parameters of cls rd, rn.
type ClsParams struct {
	Rd, Rn arm64.Reg
}

func NewClsParams(rd arm64.Reg, rn arm64.Reg) ClsParams {
	return ClsParams{
		Rd: rd,
		Rn: rn,
	}
}

func (p ClsParams) Instr() arm64.Instr {
	in, err := arm64.New().Cls(p.Rd, p.Rn)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

func (p ClsParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// clsGen — generator for cls: registers of the same width, 31st is zr.
type clsGen struct {
	rnd *rand.Rand
}

// Cls — an arbitrary cls.
func Cls(rnd *rand.Rand) ohsnap.Arbitrary[ClsParams] {
	return newClsGen(rnd)
}

func newClsGen(rnd *rand.Rand) clsGen {
	return clsGen{rnd: rnd}
}

func (g clsGen) Generate() iter.Seq[ClsParams] {
	return arb.Stream(func() ClsParams {
		is64 := g.rnd.IntN(2) == 1
		return NewClsParams(genReg(g.rnd, is64, false, true), genReg(g.rnd, is64, false, true))
	})
}

func (g clsGen) Shrink(p ClsParams) iter.Seq[ClsParams] {
	regs := regShrunk(p.Rd)
	out := make([]ClsParams, 0, len(regs))
	for _, r := range regs {
		out = append(out, NewClsParams(r, p.Rn))
	}

	for _, r := range regShrunk(p.Rn) {
		out = append(out, NewClsParams(p.Rd, r))
	}

	return slices.Values(out)
}
