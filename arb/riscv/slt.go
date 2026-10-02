package riscv

// Generator for Slt - one generator, one type, one constructor
// (over the shared three-register core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// Slt - an arbitrary Slt.
func Slt(rnd *rand.Rand) ohsnap.Arbitrary[RrrParams] {
	return Rrr(rnd, RrrSlt)
}
