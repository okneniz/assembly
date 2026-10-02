package riscv

// Generator for Mulhu - one generator, one type, one constructor
// (over the shared three-register core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// Mulhu - an arbitrary Mulhu.
func Mulhu(rnd *rand.Rand) ohsnap.Arbitrary[RrrParams] {
	return Rrr(rnd, RrrMulhu)
}
