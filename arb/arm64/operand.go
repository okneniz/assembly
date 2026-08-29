package arm64

// Operand generators for arm64: one per operand type (arch/arm64/operand.go).
// Values are always valid; the shrink is semantic — toward x0/w0, toward zero,
// toward the first constant. Dependent parts (width, classes) are decided by
// the instruction generator — the operand generators here are context-free.

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
	"github.com/okneniz/assembly/arch/arm64"
)

// --- types ---------------------------------------------------------------------

// immArb — immediate generator: uniform in [0, max], shrink by halving toward 0.
// make is a checked constructor (NewImm12/NewImm16/...): the error channel
// is propagated to Generate/Shrink.
type immArb[T any] struct {
	rnd  *rand.Rand
	max  int64
	make func(int64) (T, error)
}

// regArb — an arbitrary register of any width and class.
type regArb struct {
	rnd *rand.Rand
}

func newRegArb(rnd *rand.Rand) regArb {
	return regArb{rnd: rnd}
}

// offArb — an arbitrary offset, valid for both widths.
type offArb struct {
	rnd *rand.Rand
}

func newOffArb(rnd *rand.Rand) offArb {
	return offArb{rnd: rnd}
}

// --- constructors --------------------------------------------------------------

// Reg — an arbitrary register of any width and class (including the named
// 31st ones); positions with instruction-specific constraints are filtered
// by the instruction generator.
func Reg(rnd *rand.Rand) ohsnap.Arbitrary[arm64.Reg] {
	return newRegArb(rnd)
}

// Imm12 — an arbitrary imm12 (0..4095).
func Imm12(rnd *rand.Rand) ohsnap.Arbitrary[arm64.Imm12] {
	return immArb[arm64.Imm12]{
		rnd:  rnd,
		max:  0xfff,
		make: arm64.NewImm12,
	}
}

// Imm16 — an arbitrary imm16 (0..65535).
func Imm16(rnd *rand.Rand) ohsnap.Arbitrary[arm64.Imm16] {
	return immArb[arm64.Imm16]{
		rnd:  rnd,
		max:  0xffff,
		make: arm64.NewImm16,
	}
}

// Imm6 — an arbitrary shift amount (0..63).
func Imm6(rnd *rand.Rand) ohsnap.Arbitrary[arm64.Imm6] {
	return immArb[arm64.Imm6]{
		rnd:  rnd,
		max:  63,
		make: arm64.NewImm6,
	}
}

// Hw — an arbitrary halfword position of movz/movk.
func Hw(rnd *rand.Rand) ohsnap.Arbitrary[arm64.Hw] {
	return arb.Enum(rnd, arm64.Hw0, arm64.Hw1, arm64.Hw2, arm64.Hw3)
}

// Shift — an arbitrary shift kind.
func Shift(rnd *rand.Rand) ohsnap.Arbitrary[arm64.Shift] {
	return arb.Enum(rnd, arm64.LSL, arm64.LSR, arm64.ASR, arm64.ROR)
}

// Sh12 — an arbitrary add/sub immediate shift.
func Sh12(rnd *rand.Rand) ohsnap.Arbitrary[arm64.Sh12] {
	return arb.Enum(rnd, arm64.NoSh12, arm64.LSL12)
}

// Off — an arbitrary aligned load/store offset: 4-byte alignment and the
// range 0..0x3ffc are valid for both widths (the x-form allows more —
// 0..0x7ff8; those are covered by shrinking and by unit tests of the
// constructors).
func Off(rnd *rand.Rand) ohsnap.Arbitrary[arm64.Off] {
	return newOffArb(rnd)
}

// --- methods -------------------------------------------------------------------

func (a immArb[T]) Generate() iter.Seq[T] {
	return arb.Stream(func() T {
		v, err := a.make(a.rnd.Int64N(a.max + 1))
		if err != nil {
			var zero T // unreachable: the range [0, max] is set by the generator's constructor
			return zero
		}

		return v
	})
}

func (a immArb[T]) Shrink(v T) iter.Seq[T] {
	n, err := immValue(v)
	if err != nil {
		return ohsnap.Empty[T]() // String() of our own type is unparseable — invariant
	}

	var out []T
	// The upper bound first (decoder bugs live at the range edges),
	// then halving toward zero; zero itself is the halving target.
	for d := range shrink.Concat(
		shrink.Boundaries[int64](0, a.max),
		shrink.Halving[int64](0),
	)(n) {
		c, err := a.make(d)
		if err != nil {
			continue // unreachable: the candidates of [0, max] are in range
		}

		out = append(out, c)
	}

	return slices.Values(out)
}

func (a regArb) Generate() iter.Seq[arm64.Reg] {
	return arb.Stream(func() arm64.Reg {
		return genReg(a.rnd, a.rnd.IntN(2) == 1, true, true)
	})
}

func (a regArb) Shrink(r arm64.Reg) iter.Seq[arm64.Reg] {
	return slices.Values(regShrunk(r))
}

func (a offArb) Generate() iter.Seq[arm64.Off] {
	return arb.Stream(func() arm64.Off {
		return arm64.Off(a.rnd.Int64N(0x1000) * 4)
	})
}

func (a offArb) Shrink(o arm64.Off) iter.Seq[arm64.Off] {
	return slices.Values(offShrunk(o, 2))
}

// --- generation and shrink helpers ---------------------------------------------

// xreg/wreg — x#/w# by number; the unreachable error branch — the input is
// bounded by construction (rnd.IntN(31), n/2 of a valid register number).
func xreg(n int) arm64.Reg {
	r, err := arm64.X(n)
	if err != nil {
		return arm64.XZR // unreachable: n is always in 0..30
	}

	return r
}

func wreg(n int) arm64.Reg {
	r, err := arm64.W(n)
	if err != nil {
		return arm64.WZR // unreachable: n is always in 0..30
	}

	return r
}

// imm12/imm16/imm6 — checked immediates for the generators' Generate();
// the unreachable error branch — the range is set by the call (Int64N with a maximum).
func imm12(v int64) arm64.Imm12 {
	x, err := arm64.NewImm12(v)
	if err != nil {
		return arm64.Imm12{} // unreachable: v is always in 0..4095
	}

	return x
}

func imm16(v int64) arm64.Imm16 {
	x, err := arm64.NewImm16(v)
	if err != nil {
		return arm64.Imm16{} // unreachable: v is always in 0..65535
	}

	return x
}

func imm6(v int64) arm64.Imm6 {
	x, err := arm64.NewImm6(v)
	if err != nil {
		return arm64.Imm6{} // unreachable: v is always in 0..63
	}

	return x
}

// genReg — a register of the given width: an x/w-number, occasionally the
// named 31st (sp or zr — as allowed in the position). Dependent part of
// generation: the width is decided by the instruction, not the operand generator.
func genReg(rnd *rand.Rand, is64 bool, spOK, zrOK bool) arm64.Reg {
	if spOK && rnd.IntN(16) == 0 {
		if is64 {
			return arm64.SP
		}

		return arm64.WSP
	}

	if zrOK && rnd.IntN(16) == 0 {
		if is64 {
			return arm64.XZR
		}

		return arm64.WZR
	}

	n := rnd.IntN(31)
	if is64 {
		return xreg(n)
	}

	return wreg(n)
}

// regShrunk — register shrink candidates: x0/w0 of the same width and the
// halved number of the same class. The width is preserved — the invariant
// "instruction registers are of the same width" is not violated.
func regShrunk(r arm64.Reg) []arm64.Reg {
	zero := wreg(0)
	if r.Is64() {
		zero = xreg(0)
	}

	var out []arm64.Reg
	if r != zero {
		out = append(out, zero)
	}

	if n, ok := plainRegNum(r); ok && n > 1 {
		var half arm64.Reg
		if r.Is64() {
			half = xreg(n / 2)
		} else {
			half = wreg(n / 2)
		}

		out = append(out, half)
	}

	return out
}

// plainRegNum — the number of a plain register (x#/w#); not for named ones.
func plainRegNum(r arm64.Reg) (int, bool) {
	s := r.String()
	if len(s) < 2 || (s[0] != 'x' && s[0] != 'w') {
		return 0, false
	}

	n, err := strconv.Atoi(s[1:])
	if err != nil {
		return 0, false
	}

	return n, true
}

// offShrunk — offset shrink candidates: halving toward 0 while preserving
// alignment (an unaligned candidate would be rejected by the constructor).
func offShrunk(off arm64.Off, scale uint32) []arm64.Off {
	align := int64(1) << scale
	d := (int64(off) / 2) / align * align
	var out []arm64.Off
	for ; d > 0; d = (d / 2) / align * align {
		out = append(out, arm64.Off(d))
	}

	if off != 0 {
		out = append(out, 0)
	}

	return out
}

// immValue — the numeric value of an immediate from its String() ("#0x.."/"#..").
// The input is produced by String() of this package's types; an error means an
// invariant violation on our side and is handled by the caller (panic is forbidden).
func immValue(v any) (int64, error) {
	s := strings.TrimPrefix(fmt.Sprintf("%v", v), "#")
	if h, ok := strings.CutPrefix(s, "0x"); ok {
		return strconv.ParseInt(h, 16, 64)
	}

	return strconv.ParseInt(s, 10, 64)
}
