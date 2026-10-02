package riscv

// The I-type ALU core (rd, rs1, imm12): andi/ori/xori/slti/sltiu — one
// parameter type, the Op field names the family. The per-family
// generators live in their own files (andi.go, ...).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/riscv"
	"github.com/okneniz/assembly/disasm"
)

// RiOp — the I-type ALU mnemonic.
type RiOp int

const (
	RiAndi RiOp = iota
	RiOri
	RiXori
	RiSlti
	RiSltiu
	riOpCount
)

// RiParams — parameters of the I-type ALU families.
type RiParams struct {
	Op      RiOp
	Rd, Rs1 riscv.Reg
	Imm     riscv.Imm12
}

func NewRiParams(op RiOp, rd riscv.Reg, rs1 riscv.Reg, imm riscv.Imm12) RiParams {
	return RiParams{
		Op:  op,
		Rd:  rd,
		Rs1: rs1,
		Imm: imm,
	}
}

func (p RiParams) Instr() riscv.Instr {
	return riCalls[p.Op](riscv.New(), p)
}

func (p RiParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

var riCalls = [riOpCount]func(riscv.Builder, RiParams) riscv.Instr{
	RiAndi:  func(b riscv.Builder, p RiParams) riscv.Instr { return b.Andi(p.Rd, p.Rs1, p.Imm) },
	RiOri:   func(b riscv.Builder, p RiParams) riscv.Instr { return b.Ori(p.Rd, p.Rs1, p.Imm) },
	RiXori:  func(b riscv.Builder, p RiParams) riscv.Instr { return b.Xori(p.Rd, p.Rs1, p.Imm) },
	RiSlti:  func(b riscv.Builder, p RiParams) riscv.Instr { return b.Slti(p.Rd, p.Rs1, p.Imm) },
	RiSltiu: func(b riscv.Builder, p RiParams) riscv.Instr { return b.Sltiu(p.Rd, p.Rs1, p.Imm) },
}

// riGen — the shared generator of the I-type group.
type riGen struct {
	rnd *rand.Rand
	op  RiOp
}

// Ri — an arbitrary I-type ALU instruction of one family (the
// per-family files bind their mnemonic).
func Ri(rnd *rand.Rand, op RiOp) ohsnap.Arbitrary[RiParams] {
	return riGen{rnd: rnd, op: op}
}

func (g riGen) Generate() iter.Seq[RiParams] {
	return arb.Stream(func() RiParams {
		return NewRiParams(g.op, reg(g.rnd), reg(g.rnd), imm12(g.rnd))
	})
}

func (g riGen) Shrink(p RiParams) iter.Seq[RiParams] {
	out := make([]RiParams, 0, 8)
	for _, r := range regShrunk(p.Rd) {
		out = append(out, NewRiParams(p.Op, r, p.Rs1, p.Imm))
	}

	for _, r := range regShrunk(p.Rs1) {
		out = append(out, NewRiParams(p.Op, p.Rd, r, p.Imm))
	}

	for _, v := range immShrunk(p.Imm, riscv.New().Imm12, si12Shrink) {
		out = append(out, NewRiParams(p.Op, p.Rd, p.Rs1, v))
	}

	return slices.Values(out)
}
