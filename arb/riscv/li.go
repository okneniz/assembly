package riscv

// Generator for the li pseudo-form - one generator, one type, one
// constructor (Li): the expansion ladder of expandLi covers the signed
// 32-bit domain (addi / lui / lui+addiw, compressed where close); the
// value classes are drawn head-on. Values outside [-2^31, 2^31) need
// the slli ladder the expansion does not build - outside the family.

import (
	"fmt"
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/riscv"
)

// LiParams — parameters of li rd, val.
type LiParams struct {
	Rd  riscv.Reg
	Val int64
}

func NewLiParams(rd riscv.Reg, val int64) LiParams {
	return LiParams{
		Rd:  rd,
		Val: val,
	}
}

func (p LiParams) String() string {
	return fmt.Sprintf("li %s, %d", p.Rd, p.Val)
}

// liGen — generator for li: the signed 32-bit domain, the ladder
// classes (c.li / addi / lui-only) sampled head-on.
type liGen struct {
	rnd *rand.Rand
}

// Li — an arbitrary li.
func Li(rnd *rand.Rand) ohsnap.Arbitrary[LiParams] {
	return newLiGen(rnd)
}

func newLiGen(rnd *rand.Rand) liGen {
	return liGen{rnd: rnd}
}

func (g liGen) Generate() iter.Seq[LiParams] {
	rd := reg(g.rnd)
	switch g.rnd.IntN(4) {
	case 0: // the c.li window
		return arb.Stream(func() LiParams {
			return NewLiParams(rd, int64(g.rnd.IntN(64)-32))
		})
	case 1: // the addi window and its edges
		if g.rnd.IntN(2) == 1 {
			return arb.Stream(func() LiParams {
				return NewLiParams(rd, int64(g.rnd.IntN(4096)-2048))
			})
		}

		return arb.Stream(func() LiParams {
			return NewLiParams(rd, []int64{2047, -2048, 2048, -2049}[g.rnd.IntN(4)])
		})
	case 2: // the lui-only window (lo == 0) and the int32 edges
		if g.rnd.IntN(2) == 1 {
			return arb.Stream(func() LiParams {
				return NewLiParams(rd, (int64(g.rnd.IntN(1<<20))-(1<<19))<<12)
			})
		}

		return arb.Stream(func() LiParams {
			return NewLiParams(rd, []int64{
				1<<31 - 1, -(1 << 31), 1<<31 - 2048, -(1 << 31) + 2047,
				1 << 12, -(1 << 12), 0xfff000, -0xfff000,
			}[g.rnd.IntN(8)])
		})
	default: // the whole signed 32-bit domain
		return arb.Stream(func() LiParams {
			return NewLiParams(rd, int64(g.rnd.Uint64()&0xffffffff)-(1<<31))
		})
	}
}

func (g liGen) Shrink(p LiParams) iter.Seq[LiParams] {
	out := make([]LiParams, 0, 10)
	for _, r := range regShrunk(p.Rd) {
		out = append(out, NewLiParams(r, p.Val))
	}

	var val = p.Val
	for d := range immShrink(-(1 << 31), 1<<31-1)(val) {
		out = append(out, NewLiParams(p.Rd, d))
	}

	return slices.Values(out)
}
