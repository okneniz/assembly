package riscv

// Generator for lui — one generator, one type, one constructor
// (riscv.NewLui).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/riscv"
	"github.com/okneniz/assembly/disasm"
)

// LuiParams — parameters of lui.
type LuiParams struct {
	Rd  riscv.Reg
	Imm riscv.Imm20
}

func NewLuiParams(rd riscv.Reg, imm riscv.Imm20) LuiParams {
	return LuiParams{
		Rd:  rd,
		Imm: imm,
	}
}

func (p LuiParams) Instr() riscv.Instr {
	return riscv.NewLui(p.Rd, p.Imm)
}
func (p LuiParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// luiGen — generator for lui: register x0..x31, U-type field 0..0xfffff.
type luiGen struct {
	rnd *rand.Rand
}

func newLuiGen(rnd *rand.Rand) luiGen {
	return luiGen{rnd: rnd}
}

// Lui — an arbitrary lui.
func Lui(rnd *rand.Rand) ohsnap.Arbitrary[LuiParams] {
	return newLuiGen(rnd)
}

func (g luiGen) Generate() iter.Seq[LuiParams] {
	return arb.Stream(func() LuiParams {
		return NewLuiParams(reg(g.rnd), imm20(g.rnd))
	})
}

func (g luiGen) Shrink(p LuiParams) iter.Seq[LuiParams] {
	rd := regShrunk(p.Rd)
	imms := immShrunk(p.Imm, riscv.NewImm20, imm20Shrink)
	out := make([]LuiParams, 0, len(rd)+len(imms))
	for _, r := range rd {
		out = append(out, NewLuiParams(r, p.Imm))
	}

	for _, v := range imms {
		out = append(out, NewLuiParams(p.Rd, v))
	}

	return slices.Values(out)
}
