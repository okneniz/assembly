package riscv

// The three-register core: every (rd, rs1, rs2) family of the base ISA,
// the M extension and the AMO group — one parameter type, the Op field
// names the family, the call table binds its Builder call. In riscv any
// register is valid in any position, so the core has no lane laws. The
// per-family generators live in their own files (and.go, amoadd_w.go,
// ...).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/riscv"
	"github.com/okneniz/assembly/disasm"
)

// RrrOp — the three-register mnemonic (the family identity of the
// parameters).
type RrrOp int

const (
	RrrAnd RrrOp = iota
	RrrOr
	RrrXor
	RrrSll
	RrrSrl
	RrrSra
	RrrSlt
	RrrSltu
	RrrAddw
	RrrSubw
	RrrSllw
	RrrSrlw
	RrrSraw
	RrrMul
	RrrMulh
	RrrMulhsu
	RrrMulhu
	RrrMulw
	RrrDiv
	RrrDivu
	RrrDivw
	RrrDivuw
	RrrRem
	RrrRemu
	RrrRemw
	RrrRemuw
	RrrAmoaddW
	RrrAmoaddD
	RrrAmoandW
	RrrAmoandD
	RrrAmomaxW
	RrrAmomaxD
	RrrAmomaxuW
	RrrAmomaxuD
	RrrAmominW
	RrrAmominD
	RrrAmominuW
	RrrAmominuD
	RrrAmoorW
	RrrAmoorD
	RrrAmoswapW
	RrrAmoswapD
	RrrAmoxorW
	RrrAmoxorD
	rrrOpCount
)

// RrrParams — parameters of the three-register families.
type RrrParams struct {
	Op           RrrOp
	Rd, Rs1, Rs2 riscv.Reg
}

func NewRrrParams(op RrrOp, rd riscv.Reg, rs1 riscv.Reg, rs2 riscv.Reg) RrrParams {
	return RrrParams{
		Op:  op,
		Rd:  rd,
		Rs1: rs1,
		Rs2: rs2,
	}
}

func (p RrrParams) Instr() riscv.Instr {
	return rrrCalls[p.Op](riscv.New(), p)
}

func (p RrrParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

var rrrCalls = [rrrOpCount]func(riscv.Builder, RrrParams) riscv.Instr{
	RrrAnd:      func(b riscv.Builder, p RrrParams) riscv.Instr { return b.And(p.Rd, p.Rs1, p.Rs2) },
	RrrOr:       func(b riscv.Builder, p RrrParams) riscv.Instr { return b.Or(p.Rd, p.Rs1, p.Rs2) },
	RrrXor:      func(b riscv.Builder, p RrrParams) riscv.Instr { return b.Xor(p.Rd, p.Rs1, p.Rs2) },
	RrrSll:      func(b riscv.Builder, p RrrParams) riscv.Instr { return b.Sll(p.Rd, p.Rs1, p.Rs2) },
	RrrSrl:      func(b riscv.Builder, p RrrParams) riscv.Instr { return b.Srl(p.Rd, p.Rs1, p.Rs2) },
	RrrSra:      func(b riscv.Builder, p RrrParams) riscv.Instr { return b.Sra(p.Rd, p.Rs1, p.Rs2) },
	RrrSlt:      func(b riscv.Builder, p RrrParams) riscv.Instr { return b.Slt(p.Rd, p.Rs1, p.Rs2) },
	RrrSltu:     func(b riscv.Builder, p RrrParams) riscv.Instr { return b.Sltu(p.Rd, p.Rs1, p.Rs2) },
	RrrAddw:     func(b riscv.Builder, p RrrParams) riscv.Instr { return b.Addw(p.Rd, p.Rs1, p.Rs2) },
	RrrSubw:     func(b riscv.Builder, p RrrParams) riscv.Instr { return b.Subw(p.Rd, p.Rs1, p.Rs2) },
	RrrSllw:     func(b riscv.Builder, p RrrParams) riscv.Instr { return b.Sllw(p.Rd, p.Rs1, p.Rs2) },
	RrrSrlw:     func(b riscv.Builder, p RrrParams) riscv.Instr { return b.Srlw(p.Rd, p.Rs1, p.Rs2) },
	RrrSraw:     func(b riscv.Builder, p RrrParams) riscv.Instr { return b.Sraw(p.Rd, p.Rs1, p.Rs2) },
	RrrMul:      func(b riscv.Builder, p RrrParams) riscv.Instr { return b.Mul(p.Rd, p.Rs1, p.Rs2) },
	RrrMulh:     func(b riscv.Builder, p RrrParams) riscv.Instr { return b.Mulh(p.Rd, p.Rs1, p.Rs2) },
	RrrMulhsu:   func(b riscv.Builder, p RrrParams) riscv.Instr { return b.Mulhsu(p.Rd, p.Rs1, p.Rs2) },
	RrrMulhu:    func(b riscv.Builder, p RrrParams) riscv.Instr { return b.Mulhu(p.Rd, p.Rs1, p.Rs2) },
	RrrMulw:     func(b riscv.Builder, p RrrParams) riscv.Instr { return b.Mulw(p.Rd, p.Rs1, p.Rs2) },
	RrrDiv:      func(b riscv.Builder, p RrrParams) riscv.Instr { return b.Div(p.Rd, p.Rs1, p.Rs2) },
	RrrDivu:     func(b riscv.Builder, p RrrParams) riscv.Instr { return b.Divu(p.Rd, p.Rs1, p.Rs2) },
	RrrDivw:     func(b riscv.Builder, p RrrParams) riscv.Instr { return b.Divw(p.Rd, p.Rs1, p.Rs2) },
	RrrDivuw:    func(b riscv.Builder, p RrrParams) riscv.Instr { return b.Divuw(p.Rd, p.Rs1, p.Rs2) },
	RrrRem:      func(b riscv.Builder, p RrrParams) riscv.Instr { return b.Rem(p.Rd, p.Rs1, p.Rs2) },
	RrrRemu:     func(b riscv.Builder, p RrrParams) riscv.Instr { return b.Remu(p.Rd, p.Rs1, p.Rs2) },
	RrrRemw:     func(b riscv.Builder, p RrrParams) riscv.Instr { return b.Remw(p.Rd, p.Rs1, p.Rs2) },
	RrrRemuw:    func(b riscv.Builder, p RrrParams) riscv.Instr { return b.Remuw(p.Rd, p.Rs1, p.Rs2) },
	RrrAmoaddW:  func(b riscv.Builder, p RrrParams) riscv.Instr { return b.AmoaddW(p.Rd, p.Rs1, p.Rs2) },
	RrrAmoaddD:  func(b riscv.Builder, p RrrParams) riscv.Instr { return b.AmoaddD(p.Rd, p.Rs1, p.Rs2) },
	RrrAmoandW:  func(b riscv.Builder, p RrrParams) riscv.Instr { return b.AmoandW(p.Rd, p.Rs1, p.Rs2) },
	RrrAmoandD:  func(b riscv.Builder, p RrrParams) riscv.Instr { return b.AmoandD(p.Rd, p.Rs1, p.Rs2) },
	RrrAmomaxW:  func(b riscv.Builder, p RrrParams) riscv.Instr { return b.AmomaxW(p.Rd, p.Rs1, p.Rs2) },
	RrrAmomaxD:  func(b riscv.Builder, p RrrParams) riscv.Instr { return b.AmomaxD(p.Rd, p.Rs1, p.Rs2) },
	RrrAmomaxuW: func(b riscv.Builder, p RrrParams) riscv.Instr { return b.AmomaxuW(p.Rd, p.Rs1, p.Rs2) },
	RrrAmomaxuD: func(b riscv.Builder, p RrrParams) riscv.Instr { return b.AmomaxuD(p.Rd, p.Rs1, p.Rs2) },
	RrrAmominW:  func(b riscv.Builder, p RrrParams) riscv.Instr { return b.AmominW(p.Rd, p.Rs1, p.Rs2) },
	RrrAmominD:  func(b riscv.Builder, p RrrParams) riscv.Instr { return b.AmominD(p.Rd, p.Rs1, p.Rs2) },
	RrrAmominuW: func(b riscv.Builder, p RrrParams) riscv.Instr { return b.AmominuW(p.Rd, p.Rs1, p.Rs2) },
	RrrAmominuD: func(b riscv.Builder, p RrrParams) riscv.Instr { return b.AmominuD(p.Rd, p.Rs1, p.Rs2) },
	RrrAmoorW:   func(b riscv.Builder, p RrrParams) riscv.Instr { return b.AmoorW(p.Rd, p.Rs1, p.Rs2) },
	RrrAmoorD:   func(b riscv.Builder, p RrrParams) riscv.Instr { return b.AmoorD(p.Rd, p.Rs1, p.Rs2) },
	RrrAmoswapW: func(b riscv.Builder, p RrrParams) riscv.Instr { return b.AmoswapW(p.Rd, p.Rs1, p.Rs2) },
	RrrAmoswapD: func(b riscv.Builder, p RrrParams) riscv.Instr { return b.AmoswapD(p.Rd, p.Rs1, p.Rs2) },
	RrrAmoxorW:  func(b riscv.Builder, p RrrParams) riscv.Instr { return b.AmoxorW(p.Rd, p.Rs1, p.Rs2) },
	RrrAmoxorD:  func(b riscv.Builder, p RrrParams) riscv.Instr { return b.AmoxorD(p.Rd, p.Rs1, p.Rs2) },
}

// rrrGen — the shared generator of the group.
type rrrGen struct {
	rnd *rand.Rand
	op  RrrOp
}

// Rrr — an arbitrary three-register instruction of one family (the
// per-family files bind their mnemonic).
func Rrr(rnd *rand.Rand, op RrrOp) ohsnap.Arbitrary[RrrParams] {
	return rrrGen{rnd: rnd, op: op}
}

func (g rrrGen) Generate() iter.Seq[RrrParams] {
	return arb.Stream(func() RrrParams {
		return NewRrrParams(g.op, reg(g.rnd), reg(g.rnd), reg(g.rnd))
	})
}

func (g rrrGen) Shrink(p RrrParams) iter.Seq[RrrParams] {
	rd, rs1, rs2 := regShrunk(p.Rd), regShrunk(p.Rs1), regShrunk(p.Rs2)
	out := make([]RrrParams, 0, len(rd)+len(rs1)+len(rs2))
	for _, r := range rd {
		out = append(out, NewRrrParams(p.Op, r, p.Rs1, p.Rs2))
	}

	for _, r := range rs1 {
		out = append(out, NewRrrParams(p.Op, p.Rd, r, p.Rs2))
	}

	for _, r := range rs2 {
		out = append(out, NewRrrParams(p.Op, p.Rd, p.Rs1, r))
	}

	return slices.Values(out)
}
