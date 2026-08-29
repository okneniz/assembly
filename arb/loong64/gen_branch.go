package loong64

// Branch families: the compare-and-branch (2R + offset), compare-with-
// zero (R + offset), unconditional (b/bl, bare offset) and jirl (R + R +
// si16 byte offset). The branch ctors take the pc-relative byte OFFSET
// itself, so the families generate word-aligned offsets in each form's
// span directly
// inside each form's reach (with slack: a list property encodes the same
// directly.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"
	"github.com/okneniz/oh-snap/shrink"

	"github.com/okneniz/assembly/arb"
	arch "github.com/okneniz/assembly/arch/loong64"
	"github.com/okneniz/assembly/disasm"
)

// branchOffset — a word-aligned byte offset in ±span (span is a
// multiple of 4 well inside the form's reach).
func branchOffset(rnd *rand.Rand, span int64) int64 {
	steps := span / 4

	return 4 * (rnd.Int64N(2*steps+1) - steps)
}

// branchShrunk — offset shrink candidates: the value halved toward
// zero (sign-preserving); the halved values that lost the 4-byte
// alignment are skipped — a misaligned offset cannot encode.
func branchShrunk(off int64) []int64 {
	var out []int64
	for d := range shrink.Halving[int64](0)(off) {
		if d%4 != 0 {
			continue // unaligned offset — skip the candidate
		}

		out = append(out, d)
	}

	return out
}

// Branch2Ctor — a compare-and-branch constructor: rj, rd, off (the
// manual order swaps the registers; the ctor parameter names carry it).
type Branch2Ctor func(rj, rd arch.Reg, target int64) arch.Instr

// branch2Entry — one constructor of the compare-and-branch table.
type branch2Entry struct {
	name string
	ctor Branch2Ctor
}

// Branch2Params — parameters of a compare-and-branch instruction.
type Branch2Params struct {
	Rj, Rd arch.Reg
	Target int64
	Ctor   Branch2Ctor
}

func NewBranch2Params(rj, rd arch.Reg, target int64, ctor Branch2Ctor) Branch2Params {
	return Branch2Params{
		Rj:     rj,
		Rd:     rd,
		Target: target,
		Ctor:   ctor,
	}
}

func (p Branch2Params) Instr() arch.Instr {
	return p.Ctor(p.Rj, p.Rd, p.Target)
}

func (p Branch2Params) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// branch2Span — the ±span of the 16-bit compare-and-branch forms
// (si16 reaches ±131068; the slack covers the drift of list properties).
const branch2Span = int64(0x1fc00)

// branch2 — the compare-and-branch family (6 ctors).
var branch2 = []branch2Entry{
	{name: "beq", ctor: arch.NewBeq},
	{name: "bne", ctor: arch.NewBne},
	{name: "blt", ctor: arch.NewBlt},
	{name: "bge", ctor: arch.NewBge},
	{name: "bltu", ctor: arch.NewBltu},
	{name: "bgeu", ctor: arch.NewBgeu},
}

// branch2Gen — generator over the compare-and-branch family.
type branch2Gen struct {
	rnd *rand.Rand
}

func newBranch2Gen(rnd *rand.Rand) branch2Gen {
	return branch2Gen{rnd: rnd}
}

// Branch2 — an arbitrary compare-and-branch instruction.
func Branch2(rnd *rand.Rand) ohsnap.Arbitrary[Branch2Params] {
	return newBranch2Gen(rnd)
}

func (g branch2Gen) Generate() iter.Seq[Branch2Params] {
	return arb.Stream(func() Branch2Params {
		e := branch2[g.rnd.IntN(len(branch2))]
		return NewBranch2Params(reg(g.rnd), reg(g.rnd), branchOffset(g.rnd, branch2Span), e.ctor)
	})
}

func (g branch2Gen) Shrink(p Branch2Params) iter.Seq[Branch2Params] {
	rj, rd := regShrunk(p.Rj), regShrunk(p.Rd)
	ts := branchShrunk(p.Target)
	out := make([]Branch2Params, 0, len(rj)+len(rd)+len(ts))
	for _, r := range rj {
		out = append(out, NewBranch2Params(r, p.Rd, p.Target, p.Ctor))
	}

	for _, r := range rd {
		out = append(out, NewBranch2Params(p.Rj, r, p.Target, p.Ctor))
	}

	for _, t := range ts {
		out = append(out, NewBranch2Params(p.Rj, p.Rd, t, p.Ctor))
	}

	return slices.Values(out)
}

// Branch1Ctor — a compare-with-zero constructor: rj, target.
type Branch1Ctor func(rj arch.Reg, target int64) arch.Instr

// branch1Entry — one constructor of the compare-with-zero table.
type branch1Entry struct {
	name string
	ctor Branch1Ctor
}

// Branch1Params — parameters of a compare-with-zero instruction.
type Branch1Params struct {
	Rj     arch.Reg
	Target int64
	Ctor   Branch1Ctor
}

func NewBranch1Params(rj arch.Reg, target int64, ctor Branch1Ctor) Branch1Params {
	return Branch1Params{
		Rj:     rj,
		Target: target,
		Ctor:   ctor,
	}
}

func (p Branch1Params) Instr() arch.Instr {
	return p.Ctor(p.Rj, p.Target)
}

func (p Branch1Params) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// branch1Span — the ±span of the 21-bit compare-with-zero forms
// (offs21 reaches ±1 MiB; 512 KiB exercises the d5k16 imm[20:16] split).
const branch1Span = int64(0x80000)

// branch1 — the compare-with-zero family (2 ctors).
var branch1 = []branch1Entry{
	{name: "beqz", ctor: arch.NewBeqz},
	{name: "bnez", ctor: arch.NewBnez},
}

// branch1Gen — generator over the compare-with-zero family.
type branch1Gen struct {
	rnd *rand.Rand
}

func newBranch1Gen(rnd *rand.Rand) branch1Gen {
	return branch1Gen{rnd: rnd}
}

// Branch1 — an arbitrary compare-with-zero instruction.
func Branch1(rnd *rand.Rand) ohsnap.Arbitrary[Branch1Params] {
	return newBranch1Gen(rnd)
}

func (g branch1Gen) Generate() iter.Seq[Branch1Params] {
	return arb.Stream(func() Branch1Params {
		e := branch1[g.rnd.IntN(len(branch1))]
		return NewBranch1Params(reg(g.rnd), branchOffset(g.rnd, branch1Span), e.ctor)
	})
}

func (g branch1Gen) Shrink(p Branch1Params) iter.Seq[Branch1Params] {
	rj := regShrunk(p.Rj)
	ts := branchShrunk(p.Target)
	out := make([]Branch1Params, 0, len(rj)+len(ts))
	for _, r := range rj {
		out = append(out, NewBranch1Params(r, p.Target, p.Ctor))
	}

	for _, t := range ts {
		out = append(out, NewBranch1Params(p.Rj, t, p.Ctor))
	}

	return slices.Values(out)
}

// JumpCtor — an unconditional constructor: bare target (b/bl).
type JumpCtor func(target int64) arch.Instr

// jumpEntry — one constructor of the unconditional table.
type jumpEntry struct {
	name string
	ctor JumpCtor
}

// JumpParams — parameters of an unconditional branch instruction.
type JumpParams struct {
	Target int64
	Ctor   JumpCtor
}

func NewJumpParams(target int64, ctor JumpCtor) JumpParams {
	return JumpParams{
		Target: target,
		Ctor:   ctor,
	}
}

func (p JumpParams) Instr() arch.Instr {
	return p.Ctor(p.Target)
}

func (p JumpParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// jumpSpan — the ±span of the 26-bit unconditional forms (offs26
// reaches ±32 MiB; 4 MiB exercises the d10k16 imm[25:16] split).
const jumpSpan = int64(0x400000)

// jump — the unconditional family (2 ctors).
var jump = []jumpEntry{
	{name: "b", ctor: arch.NewB},
	{name: "bl", ctor: arch.NewBl},
}

// jumpGen — generator over the unconditional family.
type jumpGen struct {
	rnd *rand.Rand
}

func newJumpGen(rnd *rand.Rand) jumpGen {
	return jumpGen{rnd: rnd}
}

// Jump — an arbitrary unconditional branch (b/bl) instruction.
func Jump(rnd *rand.Rand) ohsnap.Arbitrary[JumpParams] {
	return newJumpGen(rnd)
}

func (g jumpGen) Generate() iter.Seq[JumpParams] {
	return arb.Stream(func() JumpParams {
		e := jump[g.rnd.IntN(len(jump))]
		return NewJumpParams(branchOffset(g.rnd, jumpSpan), e.ctor)
	})
}

func (g jumpGen) Shrink(p JumpParams) iter.Seq[JumpParams] {
	ts := branchShrunk(p.Target)
	out := make([]JumpParams, 0, len(ts))
	for _, t := range ts {
		out = append(out, NewJumpParams(t, p.Ctor))
	}

	return slices.Values(out)
}

// jirl — the indirect-jump family (a "reg, reg, role" family over the
// si16 byte offset; the offset is register-relative, not pc-relative,
// but the role bounds it exactly like a branch target field).
var jirl = []r2RoleEntry[arch.Off16]{
	{name: "jirl", ctor: func(rd, rj arch.Reg, v arch.Off16) arch.Instr {
		return arch.NewJirl(rd, rj, v.Val())
	}},
}

// Jirl — an arbitrary jirl instruction.
func Jirl(rnd *rand.Rand) ohsnap.Arbitrary[R2RoleParams[arch.Off16]] {
	return newR2RoleGen(rnd, jirl, Off16(rnd))
}
