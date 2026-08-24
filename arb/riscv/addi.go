package riscv

// Generator for addi — one generator, one type, one constructor
// (riscv.NewAddi).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/riscv"
	"github.com/okneniz/assembly/disasm"
)

// AddiParams — parameters of addi (rs1 = zero is printed as li, imm = 0 — as mv).
type AddiParams struct {
	Rd, Rs1 riscv.Reg
	Imm     riscv.Imm12
}

func NewAddiParams(rd riscv.Reg, rs1 riscv.Reg, imm riscv.Imm12) AddiParams {
	return AddiParams{
		Rd:  rd,
		Rs1: rs1,
		Imm: imm,
	}
}

func (p AddiParams) Instr() riscv.Instr {
	return riscv.NewAddi(p.Rd, p.Rs1, p.Imm)
}
func (p AddiParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// addiGen — generator for addi: registers x0..x31, immediate -2048..2047.
type addiGen struct {
	rnd *rand.Rand
}

func newAddiGen(rnd *rand.Rand) addiGen {
	return addiGen{rnd: rnd}
}

// Addi — an arbitrary addi.
func Addi(rnd *rand.Rand) ohsnap.Arbitrary[AddiParams] {
	return newAddiGen(rnd)
}

func (g addiGen) Generate() AddiParams {
	return NewAddiParams(reg(g.rnd), reg(g.rnd), imm12(g.rnd))
}

func (g addiGen) Shrink(p AddiParams) []AddiParams {
	rd, rs1 := regShrunk(p.Rd), regShrunk(p.Rs1)
	imms := immShrunk(p.Imm, riscv.NewImm12)
	out := make([]AddiParams, 0, len(rd)+len(rs1)+len(imms))
	for _, r := range rd {
		out = append(out, NewAddiParams(r, p.Rs1, p.Imm))
	}

	for _, r := range rs1 {
		out = append(out, NewAddiParams(p.Rd, r, p.Imm))
	}

	for _, v := range imms {
		out = append(out, NewAddiParams(p.Rd, p.Rs1, v))
	}

	return out
}
