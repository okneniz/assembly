package riscv

// Generator for ld — one generator, one type, one constructor
// (riscv.NewLd).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/riscv"
	"github.com/okneniz/assembly/disasm"
)

// LdParams — parameters of ld.
type LdParams struct {
	Rd, Rs1 riscv.Reg
	Off     riscv.Off
}

func NewLdParams(rd riscv.Reg, rs1 riscv.Reg, off riscv.Off) LdParams {
	return LdParams{
		Rd:  rd,
		Rs1: rs1,
		Off: off,
	}
}

func (p LdParams) Instr() riscv.Instr {
	return riscv.NewLd(p.Rd, p.Rs1, p.Off)
}
func (p LdParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// ldGen — generator for ld: registers x0..x31, byte offset -2048..2047.
type ldGen struct {
	rnd *rand.Rand
}

func newLdGen(rnd *rand.Rand) ldGen {
	return ldGen{rnd: rnd}
}

// Ld — an arbitrary ld.
func Ld(rnd *rand.Rand) ohsnap.Arbitrary[LdParams] {
	return newLdGen(rnd)
}

func (g ldGen) Generate() iter.Seq[LdParams] {
	return arb.Stream(func() LdParams {
		return NewLdParams(reg(g.rnd), reg(g.rnd), off(g.rnd))
	})
}

func (g ldGen) Shrink(p LdParams) iter.Seq[LdParams] {
	rd, rs1 := regShrunk(p.Rd), regShrunk(p.Rs1)
	offs := immShrunk(p.Off, riscv.NewOff)
	out := make([]LdParams, 0, len(rd)+len(rs1)+len(offs))
	for _, r := range rd {
		out = append(out, NewLdParams(r, p.Rs1, p.Off))
	}

	for _, r := range rs1 {
		out = append(out, NewLdParams(p.Rd, r, p.Off))
	}

	for _, v := range offs {
		out = append(out, NewLdParams(p.Rd, p.Rs1, v))
	}

	return slices.Values(out)
}
