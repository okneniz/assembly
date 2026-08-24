// Package riscv — arbitrary test generators (oh-snap) for RISC-V instructions:
// one generator, parameter type, and shrink per instruction.
package riscv

// Generator for add — one generator, one type, one constructor
// (riscv.NewAdd).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/riscv"
	"github.com/okneniz/assembly/disasm"
)

// AddParams — parameters of add.
type AddParams struct {
	Rd, Rs1, Rs2 riscv.Reg
}

func NewAddParams(rd riscv.Reg, rs1 riscv.Reg, rs2 riscv.Reg) AddParams {
	return AddParams{
		Rd:  rd,
		Rs1: rs1,
		Rs2: rs2,
	}
}

func (p AddParams) Instr() riscv.Instr {
	return riscv.NewAdd(p.Rd, p.Rs1, p.Rs2)
}
func (p AddParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// addGen — generator for add: registers uniform in x0..x31.
type addGen struct {
	rnd *rand.Rand
}

func newAddGen(rnd *rand.Rand) addGen {
	return addGen{rnd: rnd}
}

// Add — an arbitrary add.
func Add(rnd *rand.Rand) ohsnap.Arbitrary[AddParams] {
	return newAddGen(rnd)
}

func (g addGen) Generate() AddParams {
	return NewAddParams(reg(g.rnd), reg(g.rnd), reg(g.rnd))
}

func (g addGen) Shrink(p AddParams) []AddParams {
	rd, rs1, rs2 := regShrunk(p.Rd), regShrunk(p.Rs1), regShrunk(p.Rs2)
	out := make([]AddParams, 0, len(rd)+len(rs1)+len(rs2))
	for _, r := range rd {
		out = append(out, NewAddParams(r, p.Rs1, p.Rs2))
	}

	for _, r := range rs1 {
		out = append(out, NewAddParams(p.Rd, r, p.Rs2))
	}

	for _, r := range rs2 {
		out = append(out, NewAddParams(p.Rd, p.Rs1, r))
	}

	return out
}
