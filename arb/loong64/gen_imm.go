package loong64

// Immediate families: "reg, reg, imm" (the si12/ui12/si16 ALU forms and
// the ui14 CSR exchange), "reg, imm" (the si20 lu12i/pcaddi forms, the
// ui8 ldpte, the ui14 csrrd/csrwr) and the bare ui15 codes. The role
// operand generators come from operand.go; the family adds a random ctor
// of its table plus random registers.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	arch "github.com/okneniz/assembly/arch/loong64"
	"github.com/okneniz/assembly/disasm"
)

// R2RoleCtor — a "reg, reg, role" constructor.
type R2RoleCtor[T immRole] func(a, b arch.Reg, v T) arch.Instr

// r2RoleEntry — one constructor of a "reg, reg, role" family table.
type r2RoleEntry[T immRole] struct {
	name string
	ctor R2RoleCtor[T]
}

// R2RoleParams — parameters of any "reg, reg, role" instruction.
type R2RoleParams[T immRole] struct {
	A, B arch.Reg
	V    T
	Ctor R2RoleCtor[T]
}

func NewR2RoleParams[T immRole](a, b arch.Reg, v T, ctor R2RoleCtor[T]) R2RoleParams[T] {
	return R2RoleParams[T]{
		A:    a,
		B:    b,
		V:    v,
		Ctor: ctor,
	}
}

func (p R2RoleParams[T]) Instr() arch.Instr {
	return p.Ctor(p.A, p.B, p.V)
}

func (p R2RoleParams[T]) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// r2RoleGen — generator over a "reg, reg, role" family: a random ctor of
// the table, two random registers and a role value of the family's role
// generator.
type r2RoleGen[T immRole] struct {
	rnd   *rand.Rand
	ctors []r2RoleEntry[T]
	role  ohsnap.Arbitrary[T]
}

func newR2RoleGen[T immRole](
	rnd *rand.Rand,
	ctors []r2RoleEntry[T],
	role ohsnap.Arbitrary[T],
) r2RoleGen[T] {
	return r2RoleGen[T]{
		rnd:   rnd,
		ctors: ctors,
		role:  role,
	}
}

func (g r2RoleGen[T]) Generate() iter.Seq[R2RoleParams[T]] {
	return arb.Stream(func() R2RoleParams[T] {
		e := g.ctors[g.rnd.IntN(len(g.ctors))]
		return NewR2RoleParams(reg(g.rnd), reg(g.rnd), ohsnap.First(g.role.Generate()), e.ctor)
	})
}

func (g r2RoleGen[T]) Shrink(p R2RoleParams[T]) iter.Seq[R2RoleParams[T]] {
	a, b := regShrunk(p.A), regShrunk(p.B)
	vs := slices.Collect(g.role.Shrink(p.V))
	out := make([]R2RoleParams[T], 0, len(a)+len(b)+len(vs))
	for _, r := range a {
		out = append(out, NewR2RoleParams(r, p.B, p.V, p.Ctor))
	}

	for _, r := range b {
		out = append(out, NewR2RoleParams(p.A, r, p.V, p.Ctor))
	}

	for _, v := range vs {
		out = append(out, NewR2RoleParams(p.A, p.B, v, p.Ctor))
	}

	return slices.Values(out)
}

// R1RoleCtor — a "reg, role" constructor (lu12i.w and the pcaddi family,
// ldpte, csrrd/csrwr).
type R1RoleCtor[T immRole] func(r arch.Reg, v T) arch.Instr

// r1RoleEntry — one constructor of a "reg, role" family table.
type r1RoleEntry[T immRole] struct {
	name string
	ctor R1RoleCtor[T]
}

// R1RoleParams — parameters of any "reg, role" instruction.
type R1RoleParams[T immRole] struct {
	R    arch.Reg
	V    T
	Ctor R1RoleCtor[T]
}

func NewR1RoleParams[T immRole](r arch.Reg, v T, ctor R1RoleCtor[T]) R1RoleParams[T] {
	return R1RoleParams[T]{
		R:    r,
		V:    v,
		Ctor: ctor,
	}
}

func (p R1RoleParams[T]) Instr() arch.Instr {
	return p.Ctor(p.R, p.V)
}

func (p R1RoleParams[T]) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// r1RoleGen — generator over a "reg, role" family.
type r1RoleGen[T immRole] struct {
	rnd   *rand.Rand
	ctors []r1RoleEntry[T]
	role  ohsnap.Arbitrary[T]
}

func newR1RoleGen[T immRole](
	rnd *rand.Rand,
	ctors []r1RoleEntry[T],
	role ohsnap.Arbitrary[T],
) r1RoleGen[T] {
	return r1RoleGen[T]{
		rnd:   rnd,
		ctors: ctors,
		role:  role,
	}
}

func (g r1RoleGen[T]) Generate() iter.Seq[R1RoleParams[T]] {
	return arb.Stream(func() R1RoleParams[T] {
		e := g.ctors[g.rnd.IntN(len(g.ctors))]
		return NewR1RoleParams(reg(g.rnd), ohsnap.First(g.role.Generate()), e.ctor)
	})
}

func (g r1RoleGen[T]) Shrink(p R1RoleParams[T]) iter.Seq[R1RoleParams[T]] {
	rs := regShrunk(p.R)
	vs := slices.Collect(g.role.Shrink(p.V))
	out := make([]R1RoleParams[T], 0, len(rs)+len(vs))
	for _, r := range rs {
		out = append(out, NewR1RoleParams(r, p.V, p.Ctor))
	}

	for _, v := range vs {
		out = append(out, NewR1RoleParams(p.R, v, p.Ctor))
	}

	return slices.Values(out)
}

// CodeCtor — a bare ui15 code constructor (break/syscall/dbcl/dbar/ibar/
// idle).
type CodeCtor func(code arch.Code15) arch.Instr

// codeEntry — one constructor of the ui15 family table.
type codeEntry struct {
	name string
	ctor CodeCtor
}

// CodeParams — parameters of any ui15 code instruction.
type CodeParams struct {
	Code arch.Code15
	Ctor CodeCtor
}

func NewCodeParams(code arch.Code15, ctor CodeCtor) CodeParams {
	return CodeParams{
		Code: code,
		Ctor: ctor,
	}
}

func (p CodeParams) Instr() arch.Instr {
	return p.Ctor(p.Code)
}

func (p CodeParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// codeGen — generator over the ui15 family: a random ctor plus a code.
type codeGen struct {
	rnd   *rand.Rand
	ctors []codeEntry
	code  ohsnap.Arbitrary[arch.Code15]
}

func newCodeGen(rnd *rand.Rand, ctors []codeEntry) codeGen {
	return codeGen{
		rnd:   rnd,
		ctors: ctors,
		code:  Code15(rnd),
	}
}

func (g codeGen) Generate() iter.Seq[CodeParams] {
	return arb.Stream(func() CodeParams {
		e := g.ctors[g.rnd.IntN(len(g.ctors))]
		return NewCodeParams(ohsnap.First(g.code.Generate()), e.ctor)
	})
}

func (g codeGen) Shrink(p CodeParams) iter.Seq[CodeParams] {
	cs := immShrunk(p.Code, arch.New().Code15, halvingOnly)
	out := make([]CodeParams, 0, len(cs))
	for _, c := range cs {
		out = append(out, NewCodeParams(c, p.Ctor))
	}

	return slices.Values(out)
}

// aluImm12 — the si12 ALU immediate family (5 ctors).
var aluImm12 = []r2RoleEntry[arch.Imm12]{
	{name: "addi.w", ctor: arch.New().AddiW},
	{name: "addi.d", ctor: arch.New().AddiD},
	{name: "slti", ctor: arch.New().Slti},
	{name: "sltui", ctor: arch.New().Sltui},
	{name: "lu52i.d", ctor: arch.New().Lu52iD},
}

// AluImm12 — an arbitrary si12 ALU immediate instruction.
func AluImm12(rnd *rand.Rand) ohsnap.Arbitrary[R2RoleParams[arch.Imm12]] {
	return newR2RoleGen(rnd, aluImm12, Imm12(rnd))
}

// aluUImm12 — the ui12 logic immediate family (3 ctors).
var aluUImm12 = []r2RoleEntry[arch.UImm12]{
	{name: "andi", ctor: arch.New().Andi},
	{name: "ori", ctor: arch.New().Ori},
	{name: "xori", ctor: arch.New().Xori},
}

// AluUImm12 — an arbitrary ui12 logic immediate instruction.
func AluUImm12(rnd *rand.Rand) ohsnap.Arbitrary[R2RoleParams[arch.UImm12]] {
	return newR2RoleGen(rnd, aluUImm12, UImm12(rnd))
}

// aluImm16 — the si16 family (addu16i.d).
var aluImm16 = []r2RoleEntry[arch.Imm16]{
	{name: "addu16i.d", ctor: arch.New().Addu16iD},
}

// AluImm16 — an arbitrary si16 (addu16i.d) instruction.
func AluImm16(rnd *rand.Rand) ohsnap.Arbitrary[R2RoleParams[arch.Imm16]] {
	return newR2RoleGen(rnd, aluImm16, Imm16(rnd))
}

// imm20 — the si20 family: lu12i.w/lu32i.d and the pcaddi group (6 ctors).
var imm20 = []r1RoleEntry[arch.Imm20]{
	{name: "lu12i.w", ctor: arch.New().Lu12iW},
	{name: "lu32i.d", ctor: arch.New().Lu32iD},
	{name: "pcaddi", ctor: arch.New().Pcaddi},
	{name: "pcalau12i", ctor: arch.New().Pcalau12i},
	{name: "pcaddu12i", ctor: arch.New().Pcaddu12i},
	{name: "pcaddu18i", ctor: arch.New().Pcaddu18i},
}

// Imm20Instr — an arbitrary si20 (lu12i.w/lu32i.d/pcaddi family)
// instruction. The name avoids the clash with the Imm20 role generator
// of operand.go.
func Imm20Instr(rnd *rand.Rand) ohsnap.Arbitrary[R1RoleParams[arch.Imm20]] {
	return newR1RoleGen(rnd, imm20, Imm20(rnd))
}

// code15 — the ui15 code family (6 ctors; idle belongs here by shape —
// its privileged group is a decode/JSON concern, not a generator one).
var code15 = []codeEntry{
	{name: "break", ctor: arch.New().Break},
	{name: "syscall", ctor: arch.New().Syscall},
	{name: "dbcl", ctor: arch.New().Dbcl},
	{name: "dbar", ctor: arch.New().Dbar},
	{name: "ibar", ctor: arch.New().Ibar},
	{name: "idle", ctor: arch.New().Idle},
}

// Code15Instr — an arbitrary ui15 code instruction. The name avoids the
// clash with the Code15 role generator of operand.go.
func Code15Instr(rnd *rand.Rand) ohsnap.Arbitrary[CodeParams] {
	return newCodeGen(rnd, code15)
}
