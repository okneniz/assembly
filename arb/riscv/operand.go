package riscv

// Operand generators for riscv: one per operand type
// (arch/riscv/operand.go). There are no contextual constraints — in riscv
// any register is valid in any position; the shrink is semantic: X(0)/zero,
// immediates toward zero (sign-preserving).

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

// --- types ---------------------------------------------------------------------

// regArb — an arbitrary register x0..x31.
type regArb struct {
	rnd *rand.Rand
}

func newRegArb(rnd *rand.Rand) regArb {
	return regArb{rnd: rnd}
}

// immArb — signed immediate generator: uniform in [from, to].
type immArb[T any] struct {
	rnd      *rand.Rand
	from, to int64
	make     func(int64) (T, error)
}

// --- constructors --------------------------------------------------------------

// Reg — an arbitrary register x0..x31 (the ABI name comes from String()).
func Reg(rnd *rand.Rand) ohsnap.Arbitrary[riscv.Reg] {
	return newRegArb(rnd)
}

// Imm12 — an arbitrary imm12 (-2048..2047).
func Imm12(rnd *rand.Rand) ohsnap.Arbitrary[riscv.Imm12] {
	return immArb[riscv.Imm12]{
		rnd:  rnd,
		from: -2048,
		to:   2047,
		make: riscv.NewImm12,
	}
}

// Imm20 — an arbitrary U-type lui field (0..0xfffff).
func Imm20(rnd *rand.Rand) ohsnap.Arbitrary[riscv.Imm20] {
	return immArb[riscv.Imm20]{
		rnd:  rnd,
		from: 0,
		to:   0xfffff,
		make: riscv.NewImm20,
	}
}

// Off — an arbitrary load/store byte offset (-2048..2047).
func Off(rnd *rand.Rand) ohsnap.Arbitrary[riscv.Off] {
	return immArb[riscv.Off]{
		rnd:  rnd,
		from: -2048,
		to:   2047,
		make: riscv.NewOff,
	}
}

// --- methods -------------------------------------------------------------------

func (a regArb) Generate() iter.Seq[riscv.Reg] {
	return arb.Stream(func() riscv.Reg {
		return reg(a.rnd)
	})
}

func (a regArb) Shrink(r riscv.Reg) iter.Seq[riscv.Reg] {
	return slices.Values(regShrunk(r))
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
	return slices.Values(immShrunk(v, a.make))
}

// --- generation and shrink helpers ---------------------------------------------

// reg — an arbitrary register for Generate: the range is bounded by
// construction (IntN(32)), the constructor error is unreachable.
func reg(rnd *rand.Rand) riscv.Reg {
	r, err := riscv.X(rnd.IntN(32))
	if err != nil {
		return riscv.Zero // unreachable: IntN(32) is always in range
	}

	return r
}

// imm12 — an arbitrary imm12 (-2048..2047) for Generate.
func imm12(rnd *rand.Rand) riscv.Imm12 {
	v, err := riscv.NewImm12(rnd.Int64N(4096) - 2048)
	if err != nil {
		return riscv.Imm12{} // unreachable: Int64N(4096)-2048 is always in range
	}

	return v
}

// imm20 — an arbitrary U-type lui field (0..0xfffff) for Generate.
func imm20(rnd *rand.Rand) riscv.Imm20 {
	v, err := riscv.NewImm20(rnd.Int64N(0x100000))
	if err != nil {
		return riscv.Imm20{} // unreachable: Int64N(0x100000) is always in range
	}

	return v
}

// off — an arbitrary load/store byte offset (-2048..2047) for Generate.
func off(rnd *rand.Rand) riscv.Off {
	v, err := riscv.NewOff(rnd.Int64N(4096) - 2048)
	if err != nil {
		return riscv.Off{} // unreachable: Int64N(4096)-2048 is always in range
	}

	return v
}

// regShrunk — register shrink candidates: zero and the halved number.
// A candidate never equals the input — otherwise Shrink(zero) would return
// a self-loop and the oh-snap shrinker would go into infinite recursion.
func regShrunk(r riscv.Reg) []riscv.Reg {
	var out []riscv.Reg
	if r.Num() > 0 {
		out = append(out, riscv.Zero)
		if h := int(r.Num()) / 2; h > 0 {
			half, err := riscv.X(h)
			if err != nil {
				return out // unreachable: h <= 15 for a register x0..x31
			}

			out = append(out, half)
		}
	}

	return out
}

// immShrunk — halved shrink candidates for an immediate: the numeric value
// from String(), halving toward zero (shrink.Halving), re-wrapping with the
// checked constructor. A parse error is an invariant violation on our side
// (nil); a constructor error means the candidate is out of range, skip it.
func immShrunk[T any](v T, mk func(int64) (T, error)) []T {
	n, err := immValue(v)
	if err != nil {
		return nil // String() of our own type is unparseable — invariant
	}

	var out []T
	for d := range shrink.Halving[int64](0)(n) {
		t, err := mk(d)
		if err != nil {
			continue // the halved value is out of range — skip the candidate
		}

		out = append(out, t)
	}

	return out
}

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
