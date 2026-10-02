package riscv

// The memory cores: the loads (rd, rs1, off — the GPR set, the FP set)
// and the stores (rs2, rs1, off — both sets). The FP registers share
// the register number space (Reg 0..31 printed ft0/fa0/...). The
// per-family generators live in their own files (lb.go, flw.go, ...);
// sw and sd have their own generators from the first round.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/riscv"
	"github.com/okneniz/assembly/disasm"
)

// LoadOp — the load mnemonic.
type LoadOp int

const (
	LoadLb LoadOp = iota
	LoadLbu
	LoadLh
	LoadLhu
	LoadLwu
	LoadFlw
	LoadFld
	loadOpCount
)

// LoadParams — parameters of the load families.
type LoadParams struct {
	Op      LoadOp
	Rd, Rs1 riscv.Reg
	Off     riscv.Off
}

func NewLoadParams(op LoadOp, rd riscv.Reg, rs1 riscv.Reg, off riscv.Off) LoadParams {
	return LoadParams{
		Op:  op,
		Rd:  rd,
		Rs1: rs1,
		Off: off,
	}
}

func (p LoadParams) Instr() riscv.Instr {
	return loadCalls[p.Op](riscv.New(), p)
}

func (p LoadParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

var loadCalls = [loadOpCount]func(riscv.Builder, LoadParams) riscv.Instr{
	LoadLb:  func(b riscv.Builder, p LoadParams) riscv.Instr { return b.Lb(p.Rd, p.Rs1, p.Off) },
	LoadLbu: func(b riscv.Builder, p LoadParams) riscv.Instr { return b.Lbu(p.Rd, p.Rs1, p.Off) },
	LoadLh:  func(b riscv.Builder, p LoadParams) riscv.Instr { return b.Lh(p.Rd, p.Rs1, p.Off) },
	LoadLhu: func(b riscv.Builder, p LoadParams) riscv.Instr { return b.Lhu(p.Rd, p.Rs1, p.Off) },
	LoadLwu: func(b riscv.Builder, p LoadParams) riscv.Instr { return b.Lwu(p.Rd, p.Rs1, p.Off) },
	LoadFlw: func(b riscv.Builder, p LoadParams) riscv.Instr { return b.Flw(p.Rd, p.Rs1, p.Off) },
	LoadFld: func(b riscv.Builder, p LoadParams) riscv.Instr { return b.Fld(p.Rd, p.Rs1, p.Off) },
}

// loadGen — the shared generator of the load group.
type loadGen struct {
	rnd *rand.Rand
	op  LoadOp
}

// Load — an arbitrary load of one family (the per-family files bind
// their mnemonic).
func Load(rnd *rand.Rand, op LoadOp) ohsnap.Arbitrary[LoadParams] {
	return loadGen{rnd: rnd, op: op}
}

func (g loadGen) Generate() iter.Seq[LoadParams] {
	return arb.Stream(func() LoadParams {
		return NewLoadParams(g.op, reg(g.rnd), reg(g.rnd), off(g.rnd))
	})
}

func (g loadGen) Shrink(p LoadParams) iter.Seq[LoadParams] {
	out := make([]LoadParams, 0, 8)
	for _, r := range regShrunk(p.Rd) {
		out = append(out, NewLoadParams(p.Op, r, p.Rs1, p.Off))
	}

	for _, r := range regShrunk(p.Rs1) {
		out = append(out, NewLoadParams(p.Op, p.Rd, r, p.Off))
	}

	for _, v := range immShrunk(p.Off, riscv.New().Off, si12Shrink) {
		out = append(out, NewLoadParams(p.Op, p.Rd, p.Rs1, v))
	}

	return slices.Values(out)
}

// StoreOp — the store mnemonic.
type StoreOp int

const (
	StoreSb StoreOp = iota
	StoreSh
	StoreFsw
	StoreFsd
	storeOpCount
)

// StoreParams — parameters of the store families (the source register
// rides first, as in the text).
type StoreParams struct {
	Op       StoreOp
	Rs2, Rs1 riscv.Reg
	Off      riscv.Off
}

func NewStoreParams(op StoreOp, rs2 riscv.Reg, rs1 riscv.Reg, off riscv.Off) StoreParams {
	return StoreParams{
		Op:  op,
		Rs2: rs2,
		Rs1: rs1,
		Off: off,
	}
}

func (p StoreParams) Instr() riscv.Instr {
	return storeCalls[p.Op](riscv.New(), p)
}

func (p StoreParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

var storeCalls = [storeOpCount]func(riscv.Builder, StoreParams) riscv.Instr{
	StoreSb:  func(b riscv.Builder, p StoreParams) riscv.Instr { return b.Sb(p.Rs2, p.Rs1, p.Off) },
	StoreSh:  func(b riscv.Builder, p StoreParams) riscv.Instr { return b.Sh(p.Rs2, p.Rs1, p.Off) },
	StoreFsw: func(b riscv.Builder, p StoreParams) riscv.Instr { return b.Fsw(p.Rs2, p.Rs1, p.Off) },
	StoreFsd: func(b riscv.Builder, p StoreParams) riscv.Instr { return b.Fsd(p.Rs2, p.Rs1, p.Off) },
}

// storeGen — the shared generator of the store group.
type storeGen struct {
	rnd *rand.Rand
	op  StoreOp
}

// Store — an arbitrary store of one family (the per-family files bind
// their mnemonic).
func Store(rnd *rand.Rand, op StoreOp) ohsnap.Arbitrary[StoreParams] {
	return storeGen{rnd: rnd, op: op}
}

func (g storeGen) Generate() iter.Seq[StoreParams] {
	return arb.Stream(func() StoreParams {
		return NewStoreParams(g.op, reg(g.rnd), reg(g.rnd), off(g.rnd))
	})
}

func (g storeGen) Shrink(p StoreParams) iter.Seq[StoreParams] {
	out := make([]StoreParams, 0, 8)
	for _, r := range regShrunk(p.Rs2) {
		out = append(out, NewStoreParams(p.Op, r, p.Rs1, p.Off))
	}

	for _, r := range regShrunk(p.Rs1) {
		out = append(out, NewStoreParams(p.Op, p.Rs2, r, p.Off))
	}

	for _, v := range immShrunk(p.Off, riscv.New().Off, si12Shrink) {
		out = append(out, NewStoreParams(p.Op, p.Rs2, p.Rs1, v))
	}

	return slices.Values(out)
}
