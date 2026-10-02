package riscv

// The fused floating-point core (rd, rs1, rs2, rs3): fmadd/fmsub/
// fnmadd/fnmsub over .s/.d. The rounding mode is the text canon (rne,
// rm=0 - ObjDump does not print it). The per-family generators live
// in their own files (fmadd_s.go, ...).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/riscv"
	"github.com/okneniz/assembly/disasm"
)

// Fp4Op — the fused FP mnemonic.
type Fp4Op int

const (
	Fp4FmaddS Fp4Op = iota
	Fp4FmaddD
	Fp4FmsubS
	Fp4FmsubD
	Fp4FnmaddS
	Fp4FnmaddD
	Fp4FnmsubS
	Fp4FnmsubD
	fp4OpCount
)

// Fp4Params — parameters of the fused FP families.
type Fp4Params struct {
	Op                Fp4Op
	Rd, Rs1, Rs2, Rs3 riscv.Reg
}

func NewFp4Params(op Fp4Op, rd riscv.Reg, rs1 riscv.Reg, rs2 riscv.Reg, rs3 riscv.Reg) Fp4Params {
	return Fp4Params{
		Op:  op,
		Rd:  rd,
		Rs1: rs1,
		Rs2: rs2,
		Rs3: rs3,
	}
}

func (p Fp4Params) Instr() riscv.Instr {
	return fp4Calls[p.Op](riscv.New(), p)
}

func (p Fp4Params) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

var fp4Calls = [fp4OpCount]func(riscv.Builder, Fp4Params) riscv.Instr{
	Fp4FmaddS:  func(b riscv.Builder, p Fp4Params) riscv.Instr { return b.FmaddS(p.Rd, p.Rs1, p.Rs2, p.Rs3, 7) },
	Fp4FmaddD:  func(b riscv.Builder, p Fp4Params) riscv.Instr { return b.FmaddD(p.Rd, p.Rs1, p.Rs2, p.Rs3, 7) },
	Fp4FmsubS:  func(b riscv.Builder, p Fp4Params) riscv.Instr { return b.FmsubS(p.Rd, p.Rs1, p.Rs2, p.Rs3, 7) },
	Fp4FmsubD:  func(b riscv.Builder, p Fp4Params) riscv.Instr { return b.FmsubD(p.Rd, p.Rs1, p.Rs2, p.Rs3, 7) },
	Fp4FnmaddS: func(b riscv.Builder, p Fp4Params) riscv.Instr { return b.FnmaddS(p.Rd, p.Rs1, p.Rs2, p.Rs3, 7) },
	Fp4FnmaddD: func(b riscv.Builder, p Fp4Params) riscv.Instr { return b.FnmaddD(p.Rd, p.Rs1, p.Rs2, p.Rs3, 7) },
	Fp4FnmsubS: func(b riscv.Builder, p Fp4Params) riscv.Instr { return b.FnmsubS(p.Rd, p.Rs1, p.Rs2, p.Rs3, 7) },
	Fp4FnmsubD: func(b riscv.Builder, p Fp4Params) riscv.Instr { return b.FnmsubD(p.Rd, p.Rs1, p.Rs2, p.Rs3, 7) },
}

// fp4Gen — the shared generator of the fused FP group.
type fp4Gen struct {
	rnd *rand.Rand
	op  Fp4Op
}

// Fp4 — an arbitrary fused FP instruction of one family (the
// per-family files bind their mnemonic).
func Fp4(rnd *rand.Rand, op Fp4Op) ohsnap.Arbitrary[Fp4Params] {
	return fp4Gen{rnd: rnd, op: op}
}

func (g fp4Gen) Generate() iter.Seq[Fp4Params] {
	return arb.Stream(func() Fp4Params {
		return NewFp4Params(g.op, reg(g.rnd), reg(g.rnd), reg(g.rnd), reg(g.rnd))
	})
}

func (g fp4Gen) Shrink(p Fp4Params) iter.Seq[Fp4Params] {
	rd, rs1, rs2, rs3 := regShrunk(p.Rd), regShrunk(p.Rs1), regShrunk(p.Rs2), regShrunk(p.Rs3)
	out := make([]Fp4Params, 0, len(rd)+len(rs1)+len(rs2)+len(rs3))
	for _, r := range rd {
		out = append(out, NewFp4Params(p.Op, r, p.Rs1, p.Rs2, p.Rs3))
	}

	for _, r := range rs1 {
		out = append(out, NewFp4Params(p.Op, p.Rd, r, p.Rs2, p.Rs3))
	}

	for _, r := range rs2 {
		out = append(out, NewFp4Params(p.Op, p.Rd, p.Rs1, r, p.Rs3))
	}

	for _, r := range rs3 {
		out = append(out, NewFp4Params(p.Op, p.Rd, p.Rs1, p.Rs2, r))
	}

	return slices.Values(out)
}
