package riscv

// The immediate-shift core (rd, rs1, shamt): 0..63 for the X forms,
// 0..31 for the W ones — the bound is a generator invariant, the
// Builder takes the checked Imm12. The per-family generators live in
// their own files (slli.go, ...).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"
	"github.com/okneniz/oh-snap/shrink"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/riscv"
	"github.com/okneniz/assembly/disasm"
)

// ShiftOp — the immediate-shift mnemonic.
type ShiftOp int

const (
	ShiftSlli ShiftOp = iota
	ShiftSrli
	ShiftSrai
	ShiftSlliw
	ShiftSrliw
	ShiftSraiw
	shiftOpCount
)

// ShiftParams — parameters of the immediate-shift families.
type ShiftParams struct {
	Op      ShiftOp
	Rd, Rs1 riscv.Reg
	Shamt   uint32
}

func NewShiftParams(op ShiftOp, rd riscv.Reg, rs1 riscv.Reg, shamt uint32) ShiftParams {
	return ShiftParams{
		Op:    op,
		Rd:    rd,
		Rs1:   rs1,
		Shamt: shamt,
	}
}

func (p ShiftParams) Instr() riscv.Instr {
	return shiftCalls[p.Op](riscv.New(), p)
}

func (p ShiftParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// shiftImm — the shamt as the checked Imm12 the Builder takes.
func (p ShiftParams) shiftImm() riscv.Imm12 {
	v, err := riscv.New().Imm12(int64(p.Shamt))
	if err != nil {
		return riscv.Imm12{} // unreachable: Shamt <= 63
	}

	return v
}

// shiftMax — the shamt bound of a family: the X forms shift 64-bit
// registers (0..63), the W forms 32-bit ones (0..31).
func (p ShiftParams) shiftMax() uint32 {
	if p.Op == ShiftSlliw || p.Op == ShiftSrliw || p.Op == ShiftSraiw {
		return 31
	}

	return 63
}

var shiftCalls = [shiftOpCount]func(riscv.Builder, ShiftParams) riscv.Instr{
	ShiftSlli:  func(b riscv.Builder, p ShiftParams) riscv.Instr { return b.Slli(p.Rd, p.Rs1, p.shiftImm()) },
	ShiftSrli:  func(b riscv.Builder, p ShiftParams) riscv.Instr { return b.Srli(p.Rd, p.Rs1, p.shiftImm()) },
	ShiftSrai:  func(b riscv.Builder, p ShiftParams) riscv.Instr { return b.Srai(p.Rd, p.Rs1, p.shiftImm()) },
	ShiftSlliw: func(b riscv.Builder, p ShiftParams) riscv.Instr { return b.Slliw(p.Rd, p.Rs1, p.shiftImm()) },
	ShiftSrliw: func(b riscv.Builder, p ShiftParams) riscv.Instr { return b.Srliw(p.Rd, p.Rs1, p.shiftImm()) },
	ShiftSraiw: func(b riscv.Builder, p ShiftParams) riscv.Instr { return b.Sraiw(p.Rd, p.Rs1, p.shiftImm()) },
}

// shiftGen — the shared generator of the immediate-shift group.
type shiftGen struct {
	rnd *rand.Rand
	op  ShiftOp
}

// Shift — an arbitrary immediate shift of one family (the per-family
// files bind their mnemonic).
func Shift(rnd *rand.Rand, op ShiftOp) ohsnap.Arbitrary[ShiftParams] {
	return shiftGen{rnd: rnd, op: op}
}

func (g shiftGen) Generate() iter.Seq[ShiftParams] {
	return arb.Stream(func() ShiftParams {
		p := NewShiftParams(g.op, reg(g.rnd), reg(g.rnd), 0)
		return NewShiftParams(p.Op, p.Rd, p.Rs1, uint32(g.rnd.IntN(int(p.shiftMax())+1)))
	})
}

func (g shiftGen) Shrink(p ShiftParams) iter.Seq[ShiftParams] {
	out := make([]ShiftParams, 0, 8)
	for _, r := range regShrunk(p.Rd) {
		out = append(out, NewShiftParams(p.Op, r, p.Rs1, p.Shamt))
	}

	for _, r := range regShrunk(p.Rs1) {
		out = append(out, NewShiftParams(p.Op, p.Rd, r, p.Shamt))
	}

	for _, v := range shamtShrunk(p.Shamt) {
		out = append(out, NewShiftParams(p.Op, p.Rd, p.Rs1, v))
	}

	return slices.Values(out)
}

// shamtShrunk — the shamt shrink candidates: halving toward zero.
func shamtShrunk(v uint32) []uint32 {
	if v == 0 {
		return nil
	}

	var out []uint32
	for d := range shrink.Halving[int64](0)(int64(v)) {
		if d >= 0 && d <= 63 {
			out = append(out, uint32(d))
		}
	}

	return out
}
