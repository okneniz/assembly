package riscv

// Generator for sub — one generator, one type, one constructor
// (Sub).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/riscv"
	"github.com/okneniz/assembly/disasm"
)

// SubParams — parameters of sub (rs1 = zero is printed as neg).
type SubParams struct {
	Rd, Rs1, Rs2 riscv.Reg
}

func NewSubParams(rd riscv.Reg, rs1 riscv.Reg, rs2 riscv.Reg) SubParams {
	return SubParams{
		Rd:  rd,
		Rs1: rs1,
		Rs2: rs2,
	}
}

func (p SubParams) Instr() riscv.Instr {
	return riscv.New().Sub(p.Rd, p.Rs1, p.Rs2)
}
func (p SubParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// subGen — generator for sub: registers uniform in x0..x31.
type subGen struct {
	rnd *rand.Rand
}

func newSubGen(rnd *rand.Rand) subGen {
	return subGen{rnd: rnd}
}

// Sub — an arbitrary sub.
func Sub(rnd *rand.Rand) ohsnap.Arbitrary[SubParams] {
	return newSubGen(rnd)
}

func (g subGen) Generate() iter.Seq[SubParams] {
	return arb.Stream(func() SubParams {
		return NewSubParams(reg(g.rnd), reg(g.rnd), reg(g.rnd))
	})
}

func (g subGen) Shrink(p SubParams) iter.Seq[SubParams] {
	rd, rs1, rs2 := regShrunk(p.Rd), regShrunk(p.Rs1), regShrunk(p.Rs2)
	out := make([]SubParams, 0, len(rd)+len(rs1)+len(rs2))
	for _, r := range rd {
		out = append(out, NewSubParams(r, p.Rs1, p.Rs2))
	}

	for _, r := range rs1 {
		out = append(out, NewSubParams(p.Rd, r, p.Rs2))
	}

	for _, r := range rs2 {
		out = append(out, NewSubParams(p.Rd, p.Rs1, r))
	}

	return slices.Values(out)
}
