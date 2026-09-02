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
	{name: "ld.b", ctor: arch.New().LdB},
	{name: "ld.h", ctor: arch.New().LdH},
	{name: "ld.w", ctor: arch.New().LdW},
	{name: "ld.wu", ctor: arch.New().LdWu},
	{name: "ld.bu", ctor: arch.New().LdBu},
	{name: "ld.hu", ctor: arch.New().LdHu},
	{name: "ld.d", ctor: arch.New().LdD},
	{name: "st.b", ctor: arch.New().StB},
	{name: "st.h", ctor: arch.New().StH},
	{name: "st.w", ctor: arch.New().StW},
	{name: "st.d", ctor: arch.New().StD},
}

// LdSt — an arbitrary si12-offset load/store instruction.
func LdSt(rnd *rand.Rand) ohsnap.Arbitrary[R2RoleParams[arch.Imm12]] {
	return newR2RoleGen(rnd, ldSt, Imm12(rnd))
}

// ldptr — the si14-offset group: ldptr/stptr and ll/sc (8 ctors).
var ldptr = []r2RoleEntry[arch.Imm14]{
	{name: "ldptr.w", ctor: arch.New().LdptrW},
	{name: "ldptr.d", ctor: arch.New().LdptrD},
	{name: "stptr.w", ctor: arch.New().StptrW},
	{name: "stptr.d", ctor: arch.New().StptrD},
	{name: "ll.w", ctor: arch.New().LlW},
	{name: "ll.d", ctor: arch.New().LlD},
	{name: "sc.w", ctor: arch.New().ScW},
	{name: "sc.d", ctor: arch.New().ScD},
}

// Ldptr — an arbitrary si14-offset (ldptr/stptr/ll/sc) instruction.
func Ldptr(rnd *rand.Rand) ohsnap.Arbitrary[R2RoleParams[arch.Imm14]] {
	return newR2RoleGen(rnd, ldptr, Imm14(rnd))
}

// ldxStx — the 3R indexed and bounds accesses (28 ctors).
var ldxStx = []r3Entry{
	{name: "ldx.b", ctor: arch.New().LdxB},
	{name: "ldx.h", ctor: arch.New().LdxH},
	{name: "ldx.w", ctor: arch.New().LdxW},
	{name: "ldx.wu", ctor: arch.New().LdxWu},
	{name: "ldx.bu", ctor: arch.New().LdxBu},
	{name: "ldx.hu", ctor: arch.New().LdxHu},
	{name: "ldx.d", ctor: arch.New().LdxD},
	{name: "stx.b", ctor: arch.New().StxB},
	{name: "stx.h", ctor: arch.New().StxH},
	{name: "stx.w", ctor: arch.New().StxW},
	{name: "stx.d", ctor: arch.New().StxD},
	{name: "ldgt.b", ctor: arch.New().LdgtB},
	{name: "ldgt.h", ctor: arch.New().LdgtH},
	{name: "ldgt.w", ctor: arch.New().LdgtW},
	{name: "ldgt.d", ctor: arch.New().LdgtD},
	{name: "ldle.b", ctor: arch.New().LdleB},
	{name: "ldle.h", ctor: arch.New().LdleH},
	{name: "ldle.w", ctor: arch.New().LdleW},
	{name: "ldle.d", ctor: arch.New().LdleD},
	{name: "stgt.b", ctor: arch.New().StgtB},
	{name: "stgt.h", ctor: arch.New().StgtH},
	{name: "stgt.w", ctor: arch.New().StgtW},
	{name: "stgt.d", ctor: arch.New().StgtD},
	{name: "stle.b", ctor: arch.New().StleB},
	{name: "stle.h", ctor: arch.New().StleH},
	{name: "stle.w", ctor: arch.New().StleW},
	{name: "stle.d", ctor: arch.New().StleD},
}

// LdxStx — an arbitrary 3R indexed/bounds load/store instruction.
func LdxStx(rnd *rand.Rand) ohsnap.Arbitrary[R3Params] {
	return newR3Gen(rnd, ldxStx)
}

// ldAcq — the acquire/release 2R family (4 ctors).
var ldAcq = []r2Entry{
	{name: "llacq.w", ctor: arch.New().LlacqW},
	{name: "llacq.d", ctor: arch.New().LlacqD},
	{name: "screl.w", ctor: arch.New().ScrelW},
	{name: "screl.d", ctor: arch.New().ScrelD},
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
		return NewU5RI12Params(
			ohsnap.First(g.op.Generate()),
			reg(g.rnd),
			ohsnap.First(g.off.Generate()),
			e.ctor,
		)
	})
}

func (g u5ri12Gen) Shrink(p U5RI12Params) iter.Seq[U5RI12Params] {
	ops := immShrunk(p.Op, arch.New().UImm5, halvingOnly)
	rjs := regShrunk(p.Rj)
	offs := immShrunk(p.Off, arch.New().Imm12, halvingOnly)
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
	{name: "preld", ctor: arch.New().Preld},
	{name: "cacop", ctor: arch.New().Cacop},
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
	ops := immShrunk(p.Op, arch.New().UImm5, halvingOnly)
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
	{name: "preldx", ctor: arch.New().Preldx},
}

// Preldx — an arbitrary preldx instruction.
func Preldx(rnd *rand.Rand) ohsnap.Arbitrary[U5RRParams] {
	return newU5rrGen(rnd, preldx)
}
