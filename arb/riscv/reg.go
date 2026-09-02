package riscv

// Register generators: in riscv any register is valid in any position;
// the shrink is semantic: X(0)/zero, then the halved number.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/riscv"
)

// regArb — an arbitrary register x0..x31.
type regArb struct {
	rnd *rand.Rand
}

func newRegArb(rnd *rand.Rand) regArb {
	return regArb{rnd: rnd}
}

// Reg — an arbitrary register x0..x31 (the ABI name comes from String()).
func Reg(rnd *rand.Rand) ohsnap.Arbitrary[riscv.Reg] {
	return newRegArb(rnd)
}

func (a regArb) Generate() iter.Seq[riscv.Reg] {
	return arb.Stream(func() riscv.Reg {
		return reg(a.rnd)
	})
}

func (a regArb) Shrink(r riscv.Reg) iter.Seq[riscv.Reg] {
	return slices.Values(regShrunk(r))
}

// reg — an arbitrary register for Generate: the range is bounded by
// construction (IntN(32)), the constructor error is unreachable.
func reg(rnd *rand.Rand) riscv.Reg {
	r, err := riscv.X(rnd.IntN(32))
	if err != nil {
		return riscv.Zero // unreachable: IntN(32) is always in range
	}

	return r
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
