package loong64

// Immediate generators: uniform in the checked arch ranges (steps of 4
// for the word-aligned roles); the shrink is semantic — the range
// boundaries first (decoder bugs live at the range edges — sign extension,
// overflow truncation), then immediates toward zero (sign-preserving,
// re-checked through the role constructors).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"
	"github.com/okneniz/oh-snap/shrink"

	"github.com/okneniz/assembly/arb"
	arch "github.com/okneniz/assembly/arch/loong64"
)

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

// The si12 range of Imm12: shared by the generator bounds and the shrink
// boundary candidates; it mirrors the checked arch constructor Imm12.
const (
	si12Min = -2048
	si12Max = 2047
)

// Imm12 — an arbitrary si12 (-2048..2047).
func Imm12(rnd *rand.Rand) ohsnap.Arbitrary[arch.Imm12] {
	return newImmArb(rnd, si12Min, si12Max, 1, arch.New().Imm12)
}

// UImm12 — an arbitrary ui12 (0..4095).
func UImm12(rnd *rand.Rand) ohsnap.Arbitrary[arch.UImm12] {
	return newImmArb(rnd, 0, 4095, 1, arch.New().UImm12)
}

// Imm14 — an arbitrary word-scaled si14 byte offset (multiples of 4 in
// -16380..16380: ldptr/stptr/ll/sc).
func Imm14(rnd *rand.Rand) ohsnap.Arbitrary[arch.Imm14] {
	return newImmArb(rnd, -16380, 16380, 4, arch.New().Imm14)
}

// Imm16 — an arbitrary si16 (-32768..32767).
func Imm16(rnd *rand.Rand) ohsnap.Arbitrary[arch.Imm16] {
	return newImmArb(rnd, -32768, 32767, 1, arch.New().Imm16)
}

// Off16 — an arbitrary word-scaled si16 byte offset (multiples of 4 in
// -131068..131068: branches and jirl).
func Off16(rnd *rand.Rand) ohsnap.Arbitrary[arch.Off16] {
	return newImmArb(rnd, -131068, 131068, 4, arch.New().Off16)
}

// Imm20 — an arbitrary si20 (-524288..524287).
func Imm20(rnd *rand.Rand) ohsnap.Arbitrary[arch.Imm20] {
	return newImmArb(rnd, -524288, 524287, 1, arch.New().Imm20)
}

// UImm5 — an arbitrary ui5 (0..31).
func UImm5(rnd *rand.Rand) ohsnap.Arbitrary[arch.UImm5] {
	return newImmArb(rnd, 0, 31, 1, arch.New().UImm5)
}

// UImm6 — an arbitrary ui6 (0..63).
func UImm6(rnd *rand.Rand) ohsnap.Arbitrary[arch.UImm6] {
	return newImmArb(rnd, 0, 63, 1, arch.New().UImm6)
}

// UImm2 — an arbitrary ui2 (0..3).
func UImm2(rnd *rand.Rand) ohsnap.Arbitrary[arch.UImm2] {
	return newImmArb(rnd, 0, 3, 1, arch.New().UImm2)
}

// UImm3 — an arbitrary ui3 (0..7).
func UImm3(rnd *rand.Rand) ohsnap.Arbitrary[arch.UImm3] {
	return newImmArb(rnd, 0, 7, 1, arch.New().UImm3)
}

// Shift3 — an arbitrary alsl shift amount (1..4).
func Shift3(rnd *rand.Rand) ohsnap.Arbitrary[arch.Shift3] {
	return newImmArb(rnd, 1, 4, 1, arch.New().Shift3)
}

// UImm8 — an arbitrary ui8 (0..255).
func UImm8(rnd *rand.Rand) ohsnap.Arbitrary[arch.UImm8] {
	return newImmArb(rnd, 0, 255, 1, arch.New().UImm8)
}

// UImm14 — an arbitrary ui14 CSR number (0..16383).
func UImm14(rnd *rand.Rand) ohsnap.Arbitrary[arch.UImm14] {
	return newImmArb(rnd, 0, 16383, 1, arch.New().UImm14)
}

// Code15 — an arbitrary ui15 code (0..32767).
func Code15(rnd *rand.Rand) ohsnap.Arbitrary[arch.Code15] {
	return newImmArb(rnd, 0, 32767, 1, arch.New().Code15)
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
