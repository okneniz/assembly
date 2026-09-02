// Package loong64 — arbitrary test generators (oh-snap) for LoongArch
// (LA64) instructions. Unlike arb/riscv (one generator per instruction)
// the LoongArch operand shapes are uniform across whole families, so the
// generators are per SHAPE: each family holds a table of the family's
// constructors (arch/loong64 New* functions) and picks one at random on
// top of arbitrary operands (reg.go and imm.go below).
//
// Branch forms are pc-relative: the ctors take the byte offset itself,
// so the branch families generate word-aligned offsets in each form's
// span directly - Encode is pc-independent.
package loong64

// Register generators: in loong64 any register is valid in any position;
// the shrink is semantic: $zero, then the halved number.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	arch "github.com/okneniz/assembly/arch/loong64"
)

// regArb — an arbitrary register $r0..$r31.
type regArb struct {
	rnd *rand.Rand
}

func newRegArb(rnd *rand.Rand) regArb {
	return regArb{rnd: rnd}
}

// Reg — an arbitrary register $r0..$r31 (the canonical name comes from
// String()).
func Reg(rnd *rand.Rand) ohsnap.Arbitrary[arch.Reg] {
	return newRegArb(rnd)
}

func (a regArb) Generate() iter.Seq[arch.Reg] {
	return arb.Stream(func() arch.Reg {
		return reg(a.rnd)
	})
}

func (a regArb) Shrink(r arch.Reg) iter.Seq[arch.Reg] {
	return slices.Values(regShrunk(r))
}

// reg — an arbitrary register for Generate: the number is bounded by
// construction (IntN(32)), the R constructor error is unreachable.
func reg(rnd *rand.Rand) arch.Reg {
	r, err := arch.R(rnd.IntN(32))
	if err != nil {
		return arch.Zero // unreachable: IntN(32) is always in range
	}

	return r
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
			half, err := arch.R(h)
			if err != nil {
				return out // unreachable: h <= 15 for a register $r0..$r31
			}

			out = append(out, half)
		}
	}

	return out
}
