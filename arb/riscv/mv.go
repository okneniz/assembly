package riscv

// Generator for mv - one generator, one type, one constructor (Mv):
// the addi-alias pseudo-form (rd, rs2).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/riscv"
	"github.com/okneniz/assembly/disasm"
)

// MvParams — parameters of mv.
type MvParams struct {
	Rd, Rs2 riscv.Reg
}

func NewMvParams(rd riscv.Reg, rs2 riscv.Reg) MvParams {
	return MvParams{Rd: rd, Rs2: rs2}
}

func (p MvParams) Instr() riscv.Instr {
	return riscv.New().Mv(p.Rd, p.Rs2)
}

func (p MvParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Mv — an arbitrary mv.
func Mv(rnd *rand.Rand) ohsnap.Arbitrary[MvParams] {
	return mvGen{rnd: rnd}
}

type mvGen struct {
	rnd *rand.Rand
}

func (g mvGen) Generate() iter.Seq[MvParams] {
	return arb.Stream(func() MvParams {
		return NewMvParams(reg(g.rnd), reg(g.rnd))
	})
}

func (g mvGen) Shrink(p MvParams) iter.Seq[MvParams] {
	rd, rs2 := regShrunk(p.Rd), regShrunk(p.Rs2)
	out := make([]MvParams, 0, len(rd)+len(rs2))
	for _, r := range rd {
		out = append(out, NewMvParams(r, p.Rs2))
	}

	for _, r := range rs2 {
		out = append(out, NewMvParams(p.Rd, r))
	}

	return slices.Values(out)
}
