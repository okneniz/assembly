package loong64

// Bit-field families: bstrins/bstrpick (the msb >= lsb cross-check is a
// domain condition — the generator guarantees it and the shrink keeps
// it), alsl (the 1..4 shift amount) and bytepick (ui2/ui3 selectors) —
// the three-registers-plus-role shape over a random ctor.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	arch "github.com/okneniz/assembly/arch/loong64"
	"github.com/okneniz/assembly/disasm"
)

// R3RoleCtor — a three-registers-plus-role constructor.
type R3RoleCtor[T immRole] func(a, b, c arch.Reg, v T) arch.Instr

// r3RoleEntry — one constructor of a three-registers-plus-role family table.
type r3RoleEntry[T immRole] struct {
	name string
	ctor R3RoleCtor[T]
}

// R3RoleParams — parameters of any three-registers-plus-role instruction.
type R3RoleParams[T immRole] struct {
	A, B, C arch.Reg
	V       T
	Ctor    R3RoleCtor[T]
}

func NewR3RoleParams[T immRole](
	a, b, c arch.Reg,
	v T,
	ctor R3RoleCtor[T],
) R3RoleParams[T] {
	return R3RoleParams[T]{
		A:    a,
		B:    b,
		C:    c,
		V:    v,
		Ctor: ctor,
	}
}

func (p R3RoleParams[T]) Instr() arch.Instr {
	return p.Ctor(p.A, p.B, p.C, p.V)
}

func (p R3RoleParams[T]) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// r3RoleGen — generator over a three-registers-plus-role family.
type r3RoleGen[T immRole] struct {
	rnd   *rand.Rand
	ctors []r3RoleEntry[T]
	role  ohsnap.Arbitrary[T]
}

func newR3RoleGen[T immRole](
	rnd *rand.Rand,
	ctors []r3RoleEntry[T],
	role ohsnap.Arbitrary[T],
) r3RoleGen[T] {
	return r3RoleGen[T]{
		rnd:   rnd,
		ctors: ctors,
		role:  role,
	}
}

func (g r3RoleGen[T]) Generate() iter.Seq[R3RoleParams[T]] {
	return arb.Stream(func() R3RoleParams[T] {
		e := g.ctors[g.rnd.IntN(len(g.ctors))]
		return NewR3RoleParams(
			reg(g.rnd),
			reg(g.rnd),
			reg(g.rnd),
			ohsnap.First(g.role.Generate()),
			e.ctor,
		)
	})
}

func (g r3RoleGen[T]) Shrink(p R3RoleParams[T]) iter.Seq[R3RoleParams[T]] {
	a, b, c := regShrunk(p.A), regShrunk(p.B), regShrunk(p.C)
	vs := slices.Collect(g.role.Shrink(p.V))
	out := make([]R3RoleParams[T], 0, len(a)+len(b)+len(c)+len(vs))
	for _, r := range a {
		out = append(out, NewR3RoleParams(r, p.B, p.C, p.V, p.Ctor))
	}

	for _, r := range b {
		out = append(out, NewR3RoleParams(p.A, r, p.C, p.V, p.Ctor))
	}

	for _, r := range c {
		out = append(out, NewR3RoleParams(p.A, p.B, r, p.V, p.Ctor))
	}

	for _, v := range vs {
		out = append(out, NewR3RoleParams(p.A, p.B, p.C, v, p.Ctor))
	}

	return slices.Values(out)
}

// FieldWCtor — a .w bit-field constructor: rd, rj, msb, lsb.
type FieldWCtor func(rd, rj arch.Reg, msb, lsb arch.UImm5) arch.Instr

// fieldWEntry — one constructor of the .w bit-field table.
type fieldWEntry struct {
	name string
	ctor FieldWCtor
}

// FieldWParams — parameters of a .w bstrins/bstrpick instruction (the
// msb >= lsb cross-check holds for every generated and shrunk value).
type FieldWParams struct {
	Rd, Rj   arch.Reg
	Msb, Lsb arch.UImm5
	Ctor     FieldWCtor
}

func NewFieldWParams(rd, rj arch.Reg, msb, lsb arch.UImm5, ctor FieldWCtor) FieldWParams {
	return FieldWParams{
		Rd:   rd,
		Rj:   rj,
		Msb:  msb,
		Lsb:  lsb,
		Ctor: ctor,
	}
}

func (p FieldWParams) Instr() arch.Instr {
	return p.Ctor(p.Rd, p.Rj, p.Msb, p.Lsb)
}

func (p FieldWParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// fieldWGen — generator over the .w bit-field family: msb uniform in
// [lsb, 31] so the pair always satisfies the cross-check.
type fieldWGen struct {
	rnd *rand.Rand
}

func newFieldWGen(rnd *rand.Rand) fieldWGen {
	return fieldWGen{rnd: rnd}
}

// fieldW — the .w bit-field family (2 ctors).
var fieldW = []fieldWEntry{
	{name: "bstrins.w", ctor: arch.NewBstrinsW},
	{name: "bstrpick.w", ctor: arch.NewBstrpickW},
}

// FieldW — an arbitrary .w bstrins/bstrpick instruction.
func FieldW(rnd *rand.Rand) ohsnap.Arbitrary[FieldWParams] {
	return newFieldWGen(rnd)
}

func (g fieldWGen) Generate() iter.Seq[FieldWParams] {
	return arb.Stream(func() FieldWParams {
		e := fieldW[g.rnd.IntN(len(fieldW))]
		lsb := wrapImm(g.rnd.Int64N(32), arch.NewUImm5)
		msb := wrapImm(lsb.Val()+g.rnd.Int64N(32-lsb.Val()), arch.NewUImm5)
		return NewFieldWParams(reg(g.rnd), reg(g.rnd), msb, lsb, e.ctor)
	})
}

func (g fieldWGen) Shrink(p FieldWParams) iter.Seq[FieldWParams] {
	rd, rj := regShrunk(p.Rd), regShrunk(p.Rj)
	// the field gap (msb - lsb) halves toward zero: msb keeps >= lsb and
	// in range (lsb + d <= the old msb <= 31)
	var msbs []arch.UImm5
	for d := range halvingOnly(p.Msb.Val() - p.Lsb.Val()) {
		msbs = append(msbs, wrapImm(p.Lsb.Val()+d, arch.NewUImm5))
	}

	// lsb halves toward zero: the candidates stay <= msb
	lsbs := immShrunk(p.Lsb, arch.NewUImm5, halvingOnly)

	out := make([]FieldWParams, 0, len(rd)+len(rj)+len(msbs)+len(lsbs))
	for _, r := range rd {
		out = append(out, NewFieldWParams(r, p.Rj, p.Msb, p.Lsb, p.Ctor))
	}

	for _, r := range rj {
		out = append(out, NewFieldWParams(p.Rd, r, p.Msb, p.Lsb, p.Ctor))
	}

	for _, v := range msbs {
		out = append(out, NewFieldWParams(p.Rd, p.Rj, v, p.Lsb, p.Ctor))
	}

	for _, v := range lsbs {
		out = append(out, NewFieldWParams(p.Rd, p.Rj, p.Msb, v, p.Ctor))
	}

	return slices.Values(out)
}

// FieldDCtor — a .d bit-field constructor: rd, rj, msb, lsb.
type FieldDCtor func(rd, rj arch.Reg, msb, lsb arch.UImm6) arch.Instr

// fieldDEntry — one constructor of the .d bit-field table.
type fieldDEntry struct {
	name string
	ctor FieldDCtor
}

// FieldDParams — parameters of a .d bstrins/bstrpick instruction.
type FieldDParams struct {
	Rd, Rj   arch.Reg
	Msb, Lsb arch.UImm6
	Ctor     FieldDCtor
}

func NewFieldDParams(rd, rj arch.Reg, msb, lsb arch.UImm6, ctor FieldDCtor) FieldDParams {
	return FieldDParams{
		Rd:   rd,
		Rj:   rj,
		Msb:  msb,
		Lsb:  lsb,
		Ctor: ctor,
	}
}

func (p FieldDParams) Instr() arch.Instr {
	return p.Ctor(p.Rd, p.Rj, p.Msb, p.Lsb)
}

func (p FieldDParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// fieldDGen — generator over the .d bit-field family.
type fieldDGen struct {
	rnd *rand.Rand
}

func newFieldDGen(rnd *rand.Rand) fieldDGen {
	return fieldDGen{rnd: rnd}
}

// fieldD — the .d bit-field family (2 ctors).
var fieldD = []fieldDEntry{
	{name: "bstrins.d", ctor: arch.NewBstrinsD},
	{name: "bstrpick.d", ctor: arch.NewBstrpickD},
}

// FieldD — an arbitrary .d bstrins/bstrpick instruction.
func FieldD(rnd *rand.Rand) ohsnap.Arbitrary[FieldDParams] {
	return newFieldDGen(rnd)
}

func (g fieldDGen) Generate() iter.Seq[FieldDParams] {
	return arb.Stream(func() FieldDParams {
		e := fieldD[g.rnd.IntN(len(fieldD))]
		lsb := wrapImm(g.rnd.Int64N(64), arch.NewUImm6)
		msb := wrapImm(lsb.Val()+g.rnd.Int64N(64-lsb.Val()), arch.NewUImm6)
		return NewFieldDParams(reg(g.rnd), reg(g.rnd), msb, lsb, e.ctor)
	})
}

func (g fieldDGen) Shrink(p FieldDParams) iter.Seq[FieldDParams] {
	rd, rj := regShrunk(p.Rd), regShrunk(p.Rj)
	// the field gap (msb - lsb) halves toward zero: msb keeps >= lsb and
	// in range (lsb + d <= the old msb <= 63)
	var msbs []arch.UImm6
	for d := range halvingOnly(p.Msb.Val() - p.Lsb.Val()) {
		msbs = append(msbs, wrapImm(p.Lsb.Val()+d, arch.NewUImm6))
	}

	// lsb halves toward zero: the candidates stay <= msb
	lsbs := immShrunk(p.Lsb, arch.NewUImm6, halvingOnly)

	out := make([]FieldDParams, 0, len(rd)+len(rj)+len(msbs)+len(lsbs))
	for _, r := range rd {
		out = append(out, NewFieldDParams(r, p.Rj, p.Msb, p.Lsb, p.Ctor))
	}

	for _, r := range rj {
		out = append(out, NewFieldDParams(p.Rd, r, p.Msb, p.Lsb, p.Ctor))
	}

	for _, v := range msbs {
		out = append(out, NewFieldDParams(p.Rd, p.Rj, v, p.Lsb, p.Ctor))
	}

	for _, v := range lsbs {
		out = append(out, NewFieldDParams(p.Rd, p.Rj, p.Msb, v, p.Ctor))
	}

	return slices.Values(out)
}

// alsl — the alsl family over the shared three-registers-plus-role shape
// (3 ctors).
var alsl = []r3RoleEntry[arch.Shift3]{
	{name: "alsl.w", ctor: arch.NewAlslW},
	{name: "alsl.wu", ctor: arch.NewAlslWu},
	{name: "alsl.d", ctor: arch.NewAlslD},
}

// Alsl — an arbitrary alsl instruction.
func Alsl(rnd *rand.Rand) ohsnap.Arbitrary[R3RoleParams[arch.Shift3]] {
	return newR3RoleGen(rnd, alsl, Shift3(rnd))
}

// bytepickW — the bytepick.w family (1 ctor).
var bytepickW = []r3RoleEntry[arch.UImm2]{
	{name: "bytepick.w", ctor: arch.NewBytepickW},
}

// BytepickW — an arbitrary bytepick.w instruction.
func BytepickW(rnd *rand.Rand) ohsnap.Arbitrary[R3RoleParams[arch.UImm2]] {
	return newR3RoleGen(rnd, bytepickW, UImm2(rnd))
}

// bytepickD — the bytepick.d family (1 ctor).
var bytepickD = []r3RoleEntry[arch.UImm3]{
	{name: "bytepick.d", ctor: arch.NewBytepickD},
}

// BytepickD — an arbitrary bytepick.d instruction.
func BytepickD(rnd *rand.Rand) ohsnap.Arbitrary[R3RoleParams[arch.UImm3]] {
	return newR3RoleGen(rnd, bytepickD, UImm3(rnd))
}
