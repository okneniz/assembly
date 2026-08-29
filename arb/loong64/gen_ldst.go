package loong64

// Load/store families: the si12-offset loads/stores, the si14 ldptr/
// stptr/ll/sc group, the 3R indexed (ldx/stx) and bounds (ldgt/ldle/
// stgt/stle) accesses, the acquire/release 2R forms, and the memory
// hints with their ui5 leading operand (preld/cacop, preldx).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	arch "github.com/okneniz/assembly/arch/loong64"
	"github.com/okneniz/assembly/disasm"
)

// ldSt — the si12-offset load/store family (11 ctors).
var ldSt = []r2RoleEntry[arch.Imm12]{
	{name: "ld.b", ctor: arch.NewLdB},
	{name: "ld.h", ctor: arch.NewLdH},
	{name: "ld.w", ctor: arch.NewLdW},
	{name: "ld.wu", ctor: arch.NewLdWu},
	{name: "ld.bu", ctor: arch.NewLdBu},
	{name: "ld.hu", ctor: arch.NewLdHu},
	{name: "ld.d", ctor: arch.NewLdD},
	{name: "st.b", ctor: arch.NewStB},
	{name: "st.h", ctor: arch.NewStH},
	{name: "st.w", ctor: arch.NewStW},
	{name: "st.d", ctor: arch.NewStD},
}

// LdSt — an arbitrary si12-offset load/store instruction.
func LdSt(rnd *rand.Rand) ohsnap.Arbitrary[R2RoleParams[arch.Imm12]] {
	return newR2RoleGen(rnd, ldSt, Imm12(rnd))
}

// ldptr — the si14-offset group: ldptr/stptr and ll/sc (8 ctors).
var ldptr = []r2RoleEntry[arch.Imm14]{
	{name: "ldptr.w", ctor: arch.NewLdptrW},
	{name: "ldptr.d", ctor: arch.NewLdptrD},
	{name: "stptr.w", ctor: arch.NewStptrW},
	{name: "stptr.d", ctor: arch.NewStptrD},
	{name: "ll.w", ctor: arch.NewLlW},
	{name: "ll.d", ctor: arch.NewLlD},
	{name: "sc.w", ctor: arch.NewScW},
	{name: "sc.d", ctor: arch.NewScD},
}

// Ldptr — an arbitrary si14-offset (ldptr/stptr/ll/sc) instruction.
func Ldptr(rnd *rand.Rand) ohsnap.Arbitrary[R2RoleParams[arch.Imm14]] {
	return newR2RoleGen(rnd, ldptr, Imm14(rnd))
}

// ldxStx — the 3R indexed and bounds accesses (28 ctors).
var ldxStx = []r3Entry{
	{name: "ldx.b", ctor: arch.NewLdxB},
	{name: "ldx.h", ctor: arch.NewLdxH},
	{name: "ldx.w", ctor: arch.NewLdxW},
	{name: "ldx.wu", ctor: arch.NewLdxWu},
	{name: "ldx.bu", ctor: arch.NewLdxBu},
	{name: "ldx.hu", ctor: arch.NewLdxHu},
	{name: "ldx.d", ctor: arch.NewLdxD},
	{name: "stx.b", ctor: arch.NewStxB},
	{name: "stx.h", ctor: arch.NewStxH},
	{name: "stx.w", ctor: arch.NewStxW},
	{name: "stx.d", ctor: arch.NewStxD},
	{name: "ldgt.b", ctor: arch.NewLdgtB},
	{name: "ldgt.h", ctor: arch.NewLdgtH},
	{name: "ldgt.w", ctor: arch.NewLdgtW},
	{name: "ldgt.d", ctor: arch.NewLdgtD},
	{name: "ldle.b", ctor: arch.NewLdleB},
	{name: "ldle.h", ctor: arch.NewLdleH},
	{name: "ldle.w", ctor: arch.NewLdleW},
	{name: "ldle.d", ctor: arch.NewLdleD},
	{name: "stgt.b", ctor: arch.NewStgtB},
	{name: "stgt.h", ctor: arch.NewStgtH},
	{name: "stgt.w", ctor: arch.NewStgtW},
	{name: "stgt.d", ctor: arch.NewStgtD},
	{name: "stle.b", ctor: arch.NewStleB},
	{name: "stle.h", ctor: arch.NewStleH},
	{name: "stle.w", ctor: arch.NewStleW},
	{name: "stle.d", ctor: arch.NewStleD},
}

// LdxStx — an arbitrary 3R indexed/bounds load/store instruction.
func LdxStx(rnd *rand.Rand) ohsnap.Arbitrary[R3Params] {
	return newR3Gen(rnd, ldxStx)
}

// ldAcq — the acquire/release 2R family (4 ctors).
var ldAcq = []r2Entry{
	{name: "llacq.w", ctor: arch.NewLlacqW},
	{name: "llacq.d", ctor: arch.NewLlacqD},
	{name: "screl.w", ctor: arch.NewScrelW},
	{name: "screl.d", ctor: arch.NewScrelD},
}

// LdAcq — an arbitrary acquire-load/release-store instruction.
func LdAcq(rnd *rand.Rand) ohsnap.Arbitrary[R2Params] {
	return newR2Gen(rnd, ldAcq)
}

// U5RI12Ctor — a "ui5, reg, si12" constructor (preld, cacop).
type U5RI12Ctor func(op arch.UImm5, rj arch.Reg, off arch.Imm12) arch.Instr

// u5ri12Entry — one constructor of the "ui5, reg, si12" table.
type u5ri12Entry struct {
	name string
	ctor U5RI12Ctor
}

// U5RI12Params — parameters of a "ui5, reg, si12" instruction.
type U5RI12Params struct {
	Op   arch.UImm5
	Rj   arch.Reg
	Off  arch.Imm12
	Ctor U5RI12Ctor
}

func NewU5RI12Params(op arch.UImm5, rj arch.Reg, off arch.Imm12, ctor U5RI12Ctor) U5RI12Params {
	return U5RI12Params{
		Op:   op,
		Rj:   rj,
		Off:  off,
		Ctor: ctor,
	}
}

func (p U5RI12Params) Instr() arch.Instr {
	return p.Ctor(p.Op, p.Rj, p.Off)
}

func (p U5RI12Params) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// u5ri12Gen — generator over the "ui5, reg, si12" family.
type u5ri12Gen struct {
	rnd   *rand.Rand
	ctors []u5ri12Entry
	op    ohsnap.Arbitrary[arch.UImm5]
	off   ohsnap.Arbitrary[arch.Imm12]
}

func newU5RI12Gen(rnd *rand.Rand, ctors []u5ri12Entry) u5ri12Gen {
	return u5ri12Gen{
		rnd:   rnd,
		ctors: ctors,
		op:    UImm5(rnd),
		off:   Imm12(rnd),
	}
}

func (g u5ri12Gen) Generate() iter.Seq[U5RI12Params] {
	return arb.Stream(func() U5RI12Params {
		e := g.ctors[g.rnd.IntN(len(g.ctors))]
		return NewU5RI12Params(ohsnap.First(g.op.Generate()), reg(g.rnd), ohsnap.First(g.off.Generate()), e.ctor)
	})
}

func (g u5ri12Gen) Shrink(p U5RI12Params) iter.Seq[U5RI12Params] {
	ops := immShrunk(p.Op, arch.NewUImm5)
	rjs := regShrunk(p.Rj)
	offs := immShrunk(p.Off, arch.NewImm12)
	out := make([]U5RI12Params, 0, len(ops)+len(rjs)+len(offs))
	for _, v := range ops {
		out = append(out, NewU5RI12Params(v, p.Rj, p.Off, p.Ctor))
	}

	for _, r := range rjs {
		out = append(out, NewU5RI12Params(p.Op, r, p.Off, p.Ctor))
	}

	for _, v := range offs {
		out = append(out, NewU5RI12Params(p.Op, p.Rj, v, p.Ctor))
	}

	return slices.Values(out)
}

// hints — the "ui5, reg, si12" memory-hint family (2 ctors).
var hints = []u5ri12Entry{
	{name: "preld", ctor: arch.NewPreld},
	{name: "cacop", ctor: arch.NewCacop},
}

// Hints — an arbitrary preld/cacop instruction.
func Hints(rnd *rand.Rand) ohsnap.Arbitrary[U5RI12Params] {
	return newU5RI12Gen(rnd, hints)
}

// U5RRCtor — a "ui5, reg, reg" constructor (preldx, invtlb).
type U5RRCtor func(op arch.UImm5, rj, rk arch.Reg) arch.Instr

// u5rrentry — one constructor of the "ui5, reg, reg" table.
type u5rrentry struct {
	name string
	ctor U5RRCtor
}

// U5RRParams — parameters of a "ui5, reg, reg" instruction.
type U5RRParams struct {
	Op     arch.UImm5
	Rj, Rk arch.Reg
	Ctor   U5RRCtor
}

func NewU5RRParams(op arch.UImm5, rj, rk arch.Reg, ctor U5RRCtor) U5RRParams {
	return U5RRParams{
		Op:   op,
		Rj:   rj,
		Rk:   rk,
		Ctor: ctor,
	}
}

func (p U5RRParams) Instr() arch.Instr {
	return p.Ctor(p.Op, p.Rj, p.Rk)
}

func (p U5RRParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// u5rrGen — generator over the "ui5, reg, reg" family.
type u5rrGen struct {
	rnd   *rand.Rand
	ctors []u5rrentry
	op    ohsnap.Arbitrary[arch.UImm5]
}

func newU5rrGen(rnd *rand.Rand, ctors []u5rrentry) u5rrGen {
	return u5rrGen{
		rnd:   rnd,
		ctors: ctors,
		op:    UImm5(rnd),
	}
}

func (g u5rrGen) Generate() iter.Seq[U5RRParams] {
	return arb.Stream(func() U5RRParams {
		e := g.ctors[g.rnd.IntN(len(g.ctors))]
		return NewU5RRParams(ohsnap.First(g.op.Generate()), reg(g.rnd), reg(g.rnd), e.ctor)
	})
}

func (g u5rrGen) Shrink(p U5RRParams) iter.Seq[U5RRParams] {
	ops := immShrunk(p.Op, arch.NewUImm5)
	rj, rk := regShrunk(p.Rj), regShrunk(p.Rk)
	out := make([]U5RRParams, 0, len(ops)+len(rj)+len(rk))
	for _, v := range ops {
		out = append(out, NewU5RRParams(v, p.Rj, p.Rk, p.Ctor))
	}

	for _, r := range rj {
		out = append(out, NewU5RRParams(p.Op, r, p.Rk, p.Ctor))
	}

	for _, r := range rk {
		out = append(out, NewU5RRParams(p.Op, p.Rj, r, p.Ctor))
	}

	return slices.Values(out)
}

// preldx — the "ui5, reg, reg" hint family (1 ctor; invtlb of the same
// shape lives in gen_privileged.go — semantics, not shape, split them).
var preldx = []u5rrentry{
	{name: "preldx", ctor: arch.NewPreldx},
}

// Preldx — an arbitrary preldx instruction.
func Preldx(rnd *rand.Rand) ohsnap.Arbitrary[U5RRParams] {
	return newU5rrGen(rnd, preldx)
}
