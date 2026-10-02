package riscv

// Generator for Addw - one generator, one type, one constructor
// (over the shared three-register core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// Addw - an arbitrary Addw.
func Addw(rnd *rand.Rand) ohsnap.Arbitrary[RrrParams] {
	return Rrr(rnd, RrrAddw)
}
