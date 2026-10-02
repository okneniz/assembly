package riscv

// The branch core (rs1, rs2, pc-relative byte offset): even offsets in
// the 13-bit range — the bounds are generator invariants, the Builder
// stores the offset unvalidated. The per-family generators live in
// their own files (beq.go, ...).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/riscv"
	"github.com/okneniz/assembly/disasm"
)

// the pc-relative offset bounds of the branch class (even offsets: the
// RISC-V branch targets are 2-byte aligned).
const (
	branchOffMin = -4096
	branchOffMax = 4094
)

// BranchOp — the conditional branch mnemonic.
type BranchOp int

const (
	BranchBeq BranchOp = iota
	BranchBne
	BranchBlt
	BranchBge
	BranchBltu
	BranchBgeu
	branchOpCount
)

// BranchParams — parameters of the branch families.
type BranchParams struct {
	Op       BranchOp
	Rs1, Rs2 riscv.Reg
	Off      int64
}

func NewBranchParams(op BranchOp, rs1 riscv.Reg, rs2 riscv.Reg, off int64) BranchParams {
	return BranchParams{
		Op:  op,
		Rs1: rs1,
		Rs2: rs2,
		Off: off,
	}
}

func (p BranchParams) Instr() riscv.Instr {
	return branchCalls[p.Op](riscv.New(), p)
}

func (p BranchParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

var branchCalls = [branchOpCount]func(riscv.Builder, BranchParams) riscv.Instr{
	BranchBeq:  func(b riscv.Builder, p BranchParams) riscv.Instr { return b.Beq(p.Rs1, p.Rs2, p.Off) },
	BranchBne:  func(b riscv.Builder, p BranchParams) riscv.Instr { return b.Bne(p.Rs1, p.Rs2, p.Off) },
	BranchBlt:  func(b riscv.Builder, p BranchParams) riscv.Instr { return b.Blt(p.Rs1, p.Rs2, p.Off) },
	BranchBge:  func(b riscv.Builder, p BranchParams) riscv.Instr { return b.Bge(p.Rs1, p.Rs2, p.Off) },
	BranchBltu: func(b riscv.Builder, p BranchParams) riscv.Instr { return b.Bltu(p.Rs1, p.Rs2, p.Off) },
	BranchBgeu: func(b riscv.Builder, p BranchParams) riscv.Instr { return b.Bgeu(p.Rs1, p.Rs2, p.Off) },
}

// branchGen — the shared generator of the branch group.
type branchGen struct {
	rnd *rand.Rand
	op  BranchOp
}

// Branch — an arbitrary branch of one family (the per-family files
// bind their mnemonic).
func Branch(rnd *rand.Rand, op BranchOp) ohsnap.Arbitrary[BranchParams] {
	return branchGen{rnd: rnd, op: op}
}

// branchOff — an even offset in the 13-bit range.
func branchOff(rnd *rand.Rand) int64 {
	return branchOffMin + 2*rnd.Int64N((branchOffMax-branchOffMin)/2+1)
}

func (g branchGen) Generate() iter.Seq[BranchParams] {
	return arb.Stream(func() BranchParams {
		return NewBranchParams(g.op, reg(g.rnd), reg(g.rnd), branchOff(g.rnd))
	})
}

func (g branchGen) Shrink(p BranchParams) iter.Seq[BranchParams] {
	out := make([]BranchParams, 0, 8)
	for _, r := range regShrunk(p.Rs1) {
		out = append(out, NewBranchParams(p.Op, r, p.Rs2, p.Off))
	}

	for _, r := range regShrunk(p.Rs2) {
		out = append(out, NewBranchParams(p.Op, p.Rs1, r, p.Off))
	}

	for _, v := range pcOffShrunk(p.Off, branchOffMin, branchOffMax) {
		out = append(out, NewBranchParams(p.Op, p.Rs1, p.Rs2, v))
	}

	return slices.Values(out)
}

// pcOffShrunk — the pc-relative offset shrink candidates: the range
// edges, zero, then the halved value. Every candidate keeps the even
// alignment (the RISC-V targets are 2-byte aligned).
func pcOffShrunk(off int64, lo, hi int64) []int64 {
	cands := []int64{lo, hi, 0, off / 2}
	out := make([]int64, 0, len(cands))
	for _, v := range cands {
		if v%2 != 0 {
			v++
		}

		if v >= lo && v <= hi && v != off {
			out = append(out, v)
		}
	}

	return out
}
