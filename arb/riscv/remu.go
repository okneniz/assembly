package riscv

// Generator for Remu - one generator, one type, one constructor
// (over the shared three-register core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// Remu - an arbitrary Remu.
func Remu(rnd *rand.Rand) ohsnap.Arbitrary[RrrParams] {
	return Rrr(rnd, RrrRemu)
}
