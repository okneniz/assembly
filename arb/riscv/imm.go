package riscv

// Immediate generators: uniform in the checked arch ranges; the shrink is
// semantic — the range boundaries first (decoder bugs live at the range
// edges — sign extension, overflow truncation), then immediates toward
// zero (sign-preserving).

import (
	"fmt"
	"iter"
	"math/rand/v2"
	"slices"
	"strconv"
	"strings"

	ohsnap "github.com/okneniz/oh-snap"
	"github.com/okneniz/oh-snap/shrink"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/riscv"
)

// immArb — signed immediate generator: uniform in [from, to].
type immArb[T any] struct {
	rnd      *rand.Rand
	from, to int64
	make     func(int64) (T, error)
}

// Immediate ranges: the single source for the generator bounds and the
// shrink strategies; they mirror the checked arch constructors
// (Imm12/Imm20/Off).
const (
	si12Min = -2048
	si12Max = 2047
	luiMax  = 0xfffff // the U-type lui field is unsigned
)

// Imm12 — an arbitrary imm12 (-2048..2047).
func Imm12(rnd *rand.Rand) ohsnap.Arbitrary[riscv.Imm12] {
	return immArb[riscv.Imm12]{
		rnd:  rnd,
		from: si12Min,
		to:   si12Max,
		make: riscv.New().Imm12,
	}
}

// Imm20 — an arbitrary U-type lui field (0..0xfffff).
func Imm20(rnd *rand.Rand) ohsnap.Arbitrary[riscv.Imm20] {
	return immArb[riscv.Imm20]{
		rnd:  rnd,
		from: 0,
		to:   luiMax,
		make: riscv.New().Imm20,
	}
}

// Off — an arbitrary load/store byte offset (-2048..2047).
func Off(rnd *rand.Rand) ohsnap.Arbitrary[riscv.Off] {
	return immArb[riscv.Off]{
		rnd:  rnd,
		from: si12Min,
		to:   si12Max,
		make: riscv.New().Off,
	}
}

func (a immArb[T]) Generate() iter.Seq[T] {
	return arb.Stream(func() T {
		v, err := a.make(a.from + a.rnd.Int64N(a.to-a.from+1))
		if err != nil {
			var zero T // unreachable: the from..to range is bounded by construction
			return zero
		}

		return v
	})
}

func (a immArb[T]) Shrink(v T) iter.Seq[T] {
	return slices.Values(immShrunk(v, a.make, immShrink(a.from, a.to)))
}

// imm12 — an arbitrary imm12 (-2048..2047) for Generate.
func imm12(rnd *rand.Rand) riscv.Imm12 {
	v, err := riscv.New().Imm12(rnd.Int64N(4096) - 2048)
	if err != nil {
		return riscv.Imm12{} // unreachable: Int64N(4096)-2048 is always in range
	}

	return v
}

// imm20 — an arbitrary U-type lui field (0..0xfffff) for Generate.
func imm20(rnd *rand.Rand) riscv.Imm20 {
	v, err := riscv.New().Imm20(rnd.Int64N(0x100000))
	if err != nil {
		return riscv.Imm20{} // unreachable: Int64N(0x100000) is always in range
	}

	return v
}

// off — an arbitrary load/store byte offset (-2048..2047) for Generate.
func off(rnd *rand.Rand) riscv.Off {
	v, err := riscv.New().Off(rnd.Int64N(4096) - 2048)
	if err != nil {
		return riscv.Off{} // unreachable: Int64N(4096)-2048 is always in range
	}

	return v
}

// immShrunk — shrink candidates for an immediate: the strategy's numeric
// candidates, re-wrapped with the checked constructor. A parse error is an
// invariant violation on our side (nil); a constructor error means the
// candidate is out of range, skip it.
func immShrunk[T any](v T, mk func(int64) (T, error), s shrink.Shrinker[int64]) []T {
	n, err := immValue(v)
	if err != nil {
		return nil // String() of our own type is unparseable — invariant
	}

	var out []T
	for d := range s(n) {
		t, err := mk(d)
		if err != nil {
			continue // the candidate is out of range — skip it
		}

		out = append(out, t)
	}

	return out
}

// immShrink — the shrink strategy of an immediate range: the boundaries
// first (decoder bugs live at the range edges — sign extension, overflow
// truncation), then halving toward zero.
func immShrink(from, to int64) shrink.Shrinker[int64] {
	return shrink.Concat(
		shrink.Boundaries(from, to),
		shrink.Halving[int64](0),
	)
}

// Shared strategies of the instruction-level imm axes (the si12 range of
// imm12 and load/store offsets; the lui field).
var (
	si12Shrink  = immShrink(si12Min, si12Max)
	imm20Shrink = immShrink(0, luiMax)
)

// immValue — the numeric value of an immediate from its String() ("0x.."/"-0x..").
// The input is produced by String() of this package's types; an error means an
// invariant violation on our side and is handled by the caller (panic is forbidden).
func immValue(v any) (int64, error) {
	s := fmt.Sprintf("%v", v)
	neg := strings.HasPrefix(s, "-")
	s = strings.TrimPrefix(s, "-")
	s = strings.TrimPrefix(s, "0x")
	n, err := strconv.ParseInt(s, 16, 64)
	if err != nil {
		return 0, err
	}

	if neg {
		n = -n
	}

	return n, nil
}
