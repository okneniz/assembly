package riscv

// Generator for the plain jalr rs1 - one generator, one type, one
// constructor (JalrReg): the ra, 0(rs1) form (rd and the offset pinned
// by the family).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/riscv"
	"github.com/okneniz/assembly/disasm"
)

// JalrRegParams — parameters of the plain jalr rs1.
type JalrRegParams struct {
	Rs1 riscv.Reg
}

func NewJalrRegParams(rs1 riscv.Reg) JalrRegParams {
	return JalrRegParams{Rs1: rs1}
}

func (p JalrRegParams) Instr() riscv.Instr {
	return riscv.New().JalrReg(p.Rs1)
}

func (p JalrRegParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// JalrReg — an arbitrary jalr rs1.
func JalrReg(rnd *rand.Rand) ohsnap.Arbitrary[JalrRegParams] {
	return jalrRegGen{rnd: rnd}
}

type jalrRegGen struct {
	rnd *rand.Rand
}

func (g jalrRegGen) Generate() iter.Seq[JalrRegParams] {
	return arb.Stream(func() JalrRegParams {
		return NewJalrRegParams(reg(g.rnd))
	})
}

func (g jalrRegGen) Shrink(p JalrRegParams) iter.Seq[JalrRegParams] {
	out := make([]JalrRegParams, 0, 2)
	for _, r := range regShrunk(p.Rs1) {
		out = append(out, NewJalrRegParams(r))
	}

	return slices.Values(out)
}
