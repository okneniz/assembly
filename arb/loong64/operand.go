// Package loong64 — arbitrary test generators (oh-snap) for LoongArch
// (LA64) instructions. Unlike arb/riscv (one generator per instruction)
// the LoongArch operand shapes are uniform across whole families, so the
// generators are per SHAPE: each family holds a table of the family's
// constructors (arch/loong64 New* functions) and picks one at random on
// top of arbitrary operands (operand.go below).
//
// Branch forms are pc-relative: the ctors take the byte offset itself,
// so the branch families generate word-aligned offsets in each form's
// span directly - Encode is pc-independent.
package loong64

// Operand generators for loong64: one per operand role type
// (arch/loong64/operand.go). There are no contextual constraints — any
// register is valid in any position; the shrink is semantic: $zero,
// the range boundaries first, then immediates toward zero (sign-preserving,
// alignment- and range-checked through the role constructors).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"
	"github.com/okneniz/oh-snap/shrink"

	"github.com/okneniz/assembly/arb"
	arch "github.com/okneniz/assembly/arch/loong64"
)

// --- types ---------------------------------------------------------------------

// regArb — an arbitrary register $r0..$r31.
type regArb struct {
	rnd *rand.Rand
}

func newRegArb(rnd *rand.Rand) regArb {
	return regArb{rnd: rnd}
}

// immRole — the numeric value of a role operand (every arch/loong64
// immediate role exports Val).
type immRole interface {
	Val() int64
}

// immArb — a signed immediate generator: uniform in [from, to] (steps of
// `step` for the word-aligned roles), wrapped with the checked role
// constructor.
type immArb[T immRole] struct {
	rnd  *rand.Rand
	from int64
	to   int64
	step int64
	mk   func(int64) (T, error)
}

func newImmArb[T immRole](
	rnd *rand.Rand,
	from, to, step int64,
	mk func(int64) (T, error),
) immArb[T] {
	return immArb[T]{
		rnd:  rnd,
		from: from,
		to:   to,
		step: step,
		mk:   mk,
	}
}

// --- constants -------------------------------------------------------------------

// The si12 range of Imm12: shared by the generator bounds and the shrink
// boundary candidates; it mirrors the checked arch constructor NewImm12.
const (
	si12Min = -2048
	si12Max = 2047
)

// --- constructors --------------------------------------------------------------

// Reg — an arbitrary register $r0..$r31 (the canonical name comes from
// String()).
func Reg(rnd *rand.Rand) ohsnap.Arbitrary[arch.Reg] {
	return newRegArb(rnd)
}

// Imm12 — an arbitrary si12 (-2048..2047).
func Imm12(rnd *rand.Rand) ohsnap.Arbitrary[arch.Imm12] {
	return newImmArb(rnd, si12Min, si12Max, 1, arch.NewImm12)
}

// UImm12 — an arbitrary ui12 (0..4095).
func UImm12(rnd *rand.Rand) ohsnap.Arbitrary[arch.UImm12] {
	return newImmArb(rnd, 0, 4095, 1, arch.NewUImm12)
}

// Imm14 — an arbitrary word-scaled si14 byte offset (multiples of 4 in
// -16380..16380: ldptr/stptr/ll/sc).
func Imm14(rnd *rand.Rand) ohsnap.Arbitrary[arch.Imm14] {
	return newImmArb(rnd, -16380, 16380, 4, arch.NewImm14)
}

// Imm16 — an arbitrary si16 (-32768..32767).
func Imm16(rnd *rand.Rand) ohsnap.Arbitrary[arch.Imm16] {
	return newImmArb(rnd, -32768, 32767, 1, arch.NewImm16)
}

// Off16 — an arbitrary word-scaled si16 byte offset (multiples of 4 in
// -131068..131068: branches and jirl).
func Off16(rnd *rand.Rand) ohsnap.Arbitrary[arch.Off16] {
	return newImmArb(rnd, -131068, 131068, 4, arch.NewOff16)
}

// Imm20 — an arbitrary si20 (-524288..524287).
func Imm20(rnd *rand.Rand) ohsnap.Arbitrary[arch.Imm20] {
	return newImmArb(rnd, -524288, 524287, 1, arch.NewImm20)
}

// UImm5 — an arbitrary ui5 (0..31).
func UImm5(rnd *rand.Rand) ohsnap.Arbitrary[arch.UImm5] {
	return newImmArb(rnd, 0, 31, 1, arch.NewUImm5)
}

// UImm6 — an arbitrary ui6 (0..63).
func UImm6(rnd *rand.Rand) ohsnap.Arbitrary[arch.UImm6] {
	return newImmArb(rnd, 0, 63, 1, arch.NewUImm6)
}

// UImm2 — an arbitrary ui2 (0..3).
func UImm2(rnd *rand.Rand) ohsnap.Arbitrary[arch.UImm2] {
	return newImmArb(rnd, 0, 3, 1, arch.NewUImm2)
}

// UImm3 — an arbitrary ui3 (0..7).
func UImm3(rnd *rand.Rand) ohsnap.Arbitrary[arch.UImm3] {
	return newImmArb(rnd, 0, 7, 1, arch.NewUImm3)
}

// Shift3 — an arbitrary alsl shift amount (1..4).
func Shift3(rnd *rand.Rand) ohsnap.Arbitrary[arch.Shift3] {
	return newImmArb(rnd, 1, 4, 1, arch.NewShift3)
}

// UImm8 — an arbitrary ui8 (0..255).
func UImm8(rnd *rand.Rand) ohsnap.Arbitrary[arch.UImm8] {
	return newImmArb(rnd, 0, 255, 1, arch.NewUImm8)
}

// UImm14 — an arbitrary ui14 CSR number (0..16383).
func UImm14(rnd *rand.Rand) ohsnap.Arbitrary[arch.UImm14] {
	return newImmArb(rnd, 0, 16383, 1, arch.NewUImm14)
}

// Code15 — an arbitrary ui15 code (0..32767).
func Code15(rnd *rand.Rand) ohsnap.Arbitrary[arch.Code15] {
	return newImmArb(rnd, 0, 32767, 1, arch.NewCode15)
}

// --- methods -------------------------------------------------------------------

func (a regArb) Generate() iter.Seq[arch.Reg] {
	return arb.Stream(func() arch.Reg {
		return reg(a.rnd)
	})
}

func (a regArb) Shrink(r arch.Reg) iter.Seq[arch.Reg] {
	return slices.Values(regShrunk(r))
}

func (a immArb[T]) Generate() iter.Seq[T] {
	return arb.Stream(func() T {
		steps := (a.to - a.from) / a.step
		return wrapImm(a.from+a.step*a.rnd.Int64N(steps+1), a.mk)
	})
}

func (a immArb[T]) Shrink(v T) iter.Seq[T] {
	return slices.Values(immShrunk(v, a.mk, immShrink(a.from, a.to)))
}

// --- generation and shrink helpers ---------------------------------------------

// wrapImm — a role value from a number: the error of the checked
// constructor is dropped (Generate ranges are bounded by construction,
// the value is always in range). Shrink uses the checked path directly.
func wrapImm[T immRole](v int64, mk func(int64) (T, error)) T {
	x, err := mk(v) // the value is in range by construction
	if err != nil {
		panic(err) // unreachable: the Generate ranges are bounded
	}

	return x
}

// reg — an arbitrary register for Generate: the number is bounded by
// construction (IntN(32)), NewReg needs no validation.
func reg(rnd *rand.Rand) arch.Reg {
	return arch.NewReg(uint8(rnd.IntN(32)))
}

// immShrunk — shrink candidates for a role immediate: the strategy's numeric
// candidates from Val(), re-wrapped with the checked constructor. A
// constructor error (a candidate that lost alignment or left the range)
// skips the candidate.
func immShrunk[T immRole](v T, mk func(int64) (T, error), s shrink.Shrinker[int64]) []T {
	var out []T
	for d := range s(v.Val()) {
		t, err := mk(d)
		if err != nil {
			continue // the candidate is out of range — skip it
		}

		out = append(out, t)
	}

	return out
}

// --- shrink strategies ----------------------------------------------------------

// halvingOnly — the strategy of the instruction-level role axes: plain
// halving toward zero, no boundary candidates.
var halvingOnly = shrink.Halving[int64](0)

// immShrink — the shrink strategy of a standalone immediate range: the
// boundaries first (decoder bugs live at the range edges — sign extension,
// overflow truncation), then halving toward zero.
// It is used by the operand generators only: the instruction-level role
// axes stay halving-only — role boundaries ignore cross-constraints of
// the families (field msb >= lsb), a boundary candidate would break them.
func immShrink(from, to int64) shrink.Shrinker[int64] {
	return shrink.Concat(
		shrink.Boundaries(from, to),
		shrink.Halving[int64](0),
	)
}

// regShrunk — register shrink candidates: zero and the halved number.
// A candidate never equals the input — otherwise Shrink(zero) would
// return a self-loop and the oh-snap shrinker would go into infinite
// recursion.
func regShrunk(r arch.Reg) []arch.Reg {
	var out []arch.Reg
	if r.Num() > 0 {
		out = append(out, arch.Zero)
		if h := int(r.Num()) / 2; h > 0 {
			out = append(out, arch.NewReg(uint8(h)))
		}
	}

	return out
}
