package riscv

// Generator for sd — one generator, one type, one constructor
// (riscv.NewSd).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/riscv"
	"github.com/okneniz/assembly/disasm"
)

// SdParams — parameters of sd.
type SdParams struct {
	Rs2, Rs1 riscv.Reg
	Off      riscv.Off
}

func NewSdParams(rs2 riscv.Reg, rs1 riscv.Reg, off riscv.Off) SdParams {
	return SdParams{
		Rs2: rs2,
		Rs1: rs1,
		Off: off,
	}
}

func (p SdParams) Instr() riscv.Instr {
	return riscv.NewSd(p.Rs2, p.Rs1, p.Off)
}
func (p SdParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// sdGen — generator for sd: registers x0..x31, byte offset -2048..2047.
type sdGen struct {
	rnd *rand.Rand
}

func newSdGen(rnd *rand.Rand) sdGen {
	return sdGen{rnd: rnd}
}

// Sd — an arbitrary sd.
func Sd(rnd *rand.Rand) ohsnap.Arbitrary[SdParams] {
	return newSdGen(rnd)
}

func (g sdGen) Generate() iter.Seq[SdParams] {
	return arb.Stream(func() SdParams {
		return NewSdParams(reg(g.rnd), reg(g.rnd), off(g.rnd))
	})
}

func (g sdGen) Shrink(p SdParams) iter.Seq[SdParams] {
	rs2, rs1 := regShrunk(p.Rs2), regShrunk(p.Rs1)
	offs := immShrunk(p.Off, riscv.NewOff, si12Shrink)
	out := make([]SdParams, 0, len(rs2)+len(rs1)+len(offs))
	for _, r := range rs2 {
		out = append(out, NewSdParams(r, p.Rs1, p.Off))
	}

	for _, r := range rs1 {
		out = append(out, NewSdParams(p.Rs2, r, p.Off))
	}

	for _, v := range offs {
		out = append(out, NewSdParams(p.Rs2, p.Rs1, v))
	}

	return slices.Values(out)
}
