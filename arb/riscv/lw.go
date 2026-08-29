package riscv

// Generator for lw — one generator, one type, one constructor
// (riscv.NewLw).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/riscv"
	"github.com/okneniz/assembly/disasm"
)

// LwParams — parameters of lw.
type LwParams struct {
	Rd, Rs1 riscv.Reg
	Off     riscv.Off
}

func NewLwParams(rd riscv.Reg, rs1 riscv.Reg, off riscv.Off) LwParams {
	return LwParams{
		Rd:  rd,
		Rs1: rs1,
		Off: off,
	}
}

func (p LwParams) Instr() riscv.Instr {
	return riscv.NewLw(p.Rd, p.Rs1, p.Off)
}
func (p LwParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// lwGen — generator for lw: registers x0..x31, byte offset -2048..2047.
type lwGen struct {
	rnd *rand.Rand
}

func newLwGen(rnd *rand.Rand) lwGen {
	return lwGen{rnd: rnd}
}

// Lw — an arbitrary lw.
func Lw(rnd *rand.Rand) ohsnap.Arbitrary[LwParams] {
	return newLwGen(rnd)
}

func (g lwGen) Generate() iter.Seq[LwParams] {
	return arb.Stream(func() LwParams {
		return NewLwParams(reg(g.rnd), reg(g.rnd), off(g.rnd))
	})
}

func (g lwGen) Shrink(p LwParams) iter.Seq[LwParams] {
	rd, rs1 := regShrunk(p.Rd), regShrunk(p.Rs1)
	offs := immShrunk(p.Off, riscv.NewOff)
	out := make([]LwParams, 0, len(rd)+len(rs1)+len(offs))
	for _, r := range rd {
		out = append(out, NewLwParams(r, p.Rs1, p.Off))
	}

	for _, r := range rs1 {
		out = append(out, NewLwParams(p.Rd, r, p.Off))
	}

	for _, v := range offs {
		out = append(out, NewLwParams(p.Rd, p.Rs1, v))
	}

	return slices.Values(out)
}
