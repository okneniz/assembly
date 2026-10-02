package riscv

// The floating-point three-register core (rd, rs1, rs2): fadd/fsub/
// fmul/fdiv over .s/.d. The rounding mode is the text canon (rne,
// rm=0 - ObjDump does not print it). The per-family generators live
// in their own files (fadd_s.go, ...).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/riscv"
	"github.com/okneniz/assembly/disasm"
)

// Fp3Op — the three-register FP mnemonic.
type Fp3Op int

const (
	Fp3FaddS Fp3Op = iota
	Fp3FaddD
	Fp3FsubS
	Fp3FsubD
	Fp3FmulS
	Fp3FmulD
	Fp3FdivS
	Fp3FdivD
	fp3OpCount
)

// Fp3Params — parameters of the three-register FP families.
type Fp3Params struct {
	Op           Fp3Op
	Rd, Rs1, Rs2 riscv.Reg
}

func NewFp3Params(op Fp3Op, rd riscv.Reg, rs1 riscv.Reg, rs2 riscv.Reg) Fp3Params {
	return Fp3Params{
		Op:  op,
		Rd:  rd,
		Rs1: rs1,
		Rs2: rs2,
	}
}

func (p Fp3Params) Instr() riscv.Instr {
	return fp3Calls[p.Op](riscv.New(), p)
}

func (p Fp3Params) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

var fp3Calls = [fp3OpCount]func(riscv.Builder, Fp3Params) riscv.Instr{
	Fp3FaddS: func(b riscv.Builder, p Fp3Params) riscv.Instr { return b.FaddS(p.Rd, p.Rs1, p.Rs2, 7) },
	Fp3FaddD: func(b riscv.Builder, p Fp3Params) riscv.Instr { return b.FaddD(p.Rd, p.Rs1, p.Rs2, 7) },
	Fp3FsubS: func(b riscv.Builder, p Fp3Params) riscv.Instr { return b.FsubS(p.Rd, p.Rs1, p.Rs2, 7) },
	Fp3FsubD: func(b riscv.Builder, p Fp3Params) riscv.Instr { return b.FsubD(p.Rd, p.Rs1, p.Rs2, 7) },
	Fp3FmulS: func(b riscv.Builder, p Fp3Params) riscv.Instr { return b.FmulS(p.Rd, p.Rs1, p.Rs2, 7) },
	Fp3FmulD: func(b riscv.Builder, p Fp3Params) riscv.Instr { return b.FmulD(p.Rd, p.Rs1, p.Rs2, 7) },
	Fp3FdivS: func(b riscv.Builder, p Fp3Params) riscv.Instr { return b.FdivS(p.Rd, p.Rs1, p.Rs2, 7) },
	Fp3FdivD: func(b riscv.Builder, p Fp3Params) riscv.Instr { return b.FdivD(p.Rd, p.Rs1, p.Rs2, 7) },
}

// fp3Gen — the shared generator of the three-register FP group.
type fp3Gen struct {
	rnd *rand.Rand
	op  Fp3Op
}

// Fp3 — an arbitrary three-register FP instruction of one family (the
// per-family files bind their mnemonic).
func Fp3(rnd *rand.Rand, op Fp3Op) ohsnap.Arbitrary[Fp3Params] {
	return fp3Gen{rnd: rnd, op: op}
}

func (g fp3Gen) Generate() iter.Seq[Fp3Params] {
	return arb.Stream(func() Fp3Params {
		return NewFp3Params(g.op, reg(g.rnd), reg(g.rnd), reg(g.rnd))
	})
}

func (g fp3Gen) Shrink(p Fp3Params) iter.Seq[Fp3Params] {
	rd, rs1, rs2 := regShrunk(p.Rd), regShrunk(p.Rs1), regShrunk(p.Rs2)
	out := make([]Fp3Params, 0, len(rd)+len(rs1)+len(rs2))
	for _, r := range rd {
		out = append(out, NewFp3Params(p.Op, r, p.Rs1, p.Rs2))
	}

	for _, r := range rs1 {
		out = append(out, NewFp3Params(p.Op, p.Rd, r, p.Rs2))
	}

	for _, r := range rs2 {
		out = append(out, NewFp3Params(p.Op, p.Rd, p.Rs1, r))
	}

	return slices.Values(out)
}
