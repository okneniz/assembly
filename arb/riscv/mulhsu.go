package riscv

// Generator for Mulhsu - one generator, one type, one constructor
// (over the shared three-register core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// Mulhsu - an arbitrary Mulhsu.
func Mulhsu(rnd *rand.Rand) ohsnap.Arbitrary[RrrParams] {
	return Rrr(rnd, RrrMulhsu)
}
