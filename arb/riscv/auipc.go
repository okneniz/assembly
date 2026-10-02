package riscv

// Generator for auipc - one generator, one type, one constructor
// (Auipc): rd and the 20-bit U-type immediate.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/riscv"
	"github.com/okneniz/assembly/disasm"
)

// AuipcParams — parameters of auipc.
type AuipcParams struct {
	Rd  riscv.Reg
	Imm riscv.Imm20
}

func NewAuipcParams(rd riscv.Reg, imm riscv.Imm20) AuipcParams {
	return AuipcParams{Rd: rd, Imm: imm}
}

func (p AuipcParams) Instr() riscv.Instr {
	return riscv.New().Auipc(p.Rd, p.Imm)
}

func (p AuipcParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Auipc — an arbitrary auipc.
func Auipc(rnd *rand.Rand) ohsnap.Arbitrary[AuipcParams] {
	return auipcGen{rnd: rnd}
}

type auipcGen struct {
	rnd *rand.Rand
}

func (g auipcGen) Generate() iter.Seq[AuipcParams] {
	return arb.Stream(func() AuipcParams {
		return NewAuipcParams(reg(g.rnd), imm20(g.rnd))
	})
}

func (g auipcGen) Shrink(p AuipcParams) iter.Seq[AuipcParams] {
	out := make([]AuipcParams, 0, 8)
	for _, r := range regShrunk(p.Rd) {
		out = append(out, NewAuipcParams(r, p.Imm))
	}

	for _, v := range immShrunk(p.Imm, riscv.New().Imm20, imm20Shrink) {
		out = append(out, NewAuipcParams(p.Rd, v))
	}

	return slices.Values(out)
}
