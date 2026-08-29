package riscv

// Generator for sw — one generator, one type, one constructor
// (riscv.NewSw).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/riscv"
	"github.com/okneniz/assembly/disasm"
)

// SwParams — parameters of sw.
type SwParams struct {
	Rs2, Rs1 riscv.Reg
	Off      riscv.Off
}

func NewSwParams(rs2 riscv.Reg, rs1 riscv.Reg, off riscv.Off) SwParams {
	return SwParams{
		Rs2: rs2,
		Rs1: rs1,
		Off: off,
	}
}

func (p SwParams) Instr() riscv.Instr {
	return riscv.NewSw(p.Rs2, p.Rs1, p.Off)
}
func (p SwParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// swGen — generator for sw: registers x0..x31, byte offset -2048..2047.
type swGen struct {
	rnd *rand.Rand
}

func newSwGen(rnd *rand.Rand) swGen {
	return swGen{rnd: rnd}
}

// Sw — an arbitrary sw.
func Sw(rnd *rand.Rand) ohsnap.Arbitrary[SwParams] {
	return newSwGen(rnd)
}

func (g swGen) Generate() iter.Seq[SwParams] {
	return arb.Stream(func() SwParams {
		return NewSwParams(reg(g.rnd), reg(g.rnd), off(g.rnd))
	})
}

func (g swGen) Shrink(p SwParams) iter.Seq[SwParams] {
	rs2, rs1 := regShrunk(p.Rs2), regShrunk(p.Rs1)
	offs := immShrunk(p.Off, riscv.NewOff)
	out := make([]SwParams, 0, len(rs2)+len(rs1)+len(offs))
	for _, r := range rs2 {
		out = append(out, NewSwParams(r, p.Rs1, p.Off))
	}

	for _, r := range rs1 {
		out = append(out, NewSwParams(p.Rs2, r, p.Off))
	}

	for _, v := range offs {
		out = append(out, NewSwParams(p.Rs2, p.Rs1, v))
	}

	return slices.Values(out)
}
