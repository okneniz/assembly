package riscv

// Generator for jalr - one generator, one type, one constructor
// (Jalr): rd, rs1 and a 12-bit offset.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/riscv"
	"github.com/okneniz/assembly/disasm"
)

// JalrParams — parameters of jalr.
type JalrParams struct {
	Rd, Rs1 riscv.Reg
	Off     riscv.Off
}

func NewJalrParams(rd riscv.Reg, rs1 riscv.Reg, off riscv.Off) JalrParams {
	return JalrParams{
		Rd:  rd,
		Rs1: rs1,
		Off: off,
	}
}

func (p JalrParams) Instr() riscv.Instr {
	return riscv.New().Jalr(p.Rd, p.Rs1, p.Off)
}

func (p JalrParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Jalr — an arbitrary jalr.
func Jalr(rnd *rand.Rand) ohsnap.Arbitrary[JalrParams] {
	return jalrGen{rnd: rnd}
}

type jalrGen struct {
	rnd *rand.Rand
}

func (g jalrGen) Generate() iter.Seq[JalrParams] {
	return arb.Stream(func() JalrParams {
		return NewJalrParams(reg(g.rnd), reg(g.rnd), off(g.rnd))
	})
}

func (g jalrGen) Shrink(p JalrParams) iter.Seq[JalrParams] {
	out := make([]JalrParams, 0, 8)
	for _, r := range regShrunk(p.Rd) {
		out = append(out, NewJalrParams(r, p.Rs1, p.Off))
	}

	for _, r := range regShrunk(p.Rs1) {
		out = append(out, NewJalrParams(p.Rd, r, p.Off))
	}

	for _, v := range immShrunk(p.Off, riscv.New().Off, si12Shrink) {
		out = append(out, NewJalrParams(p.Rd, p.Rs1, v))
	}

	return slices.Values(out)
}
