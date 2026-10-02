package riscv

// Generator for And - one generator, one type, one constructor
// (over the shared three-register core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// And - an arbitrary And.
func And(rnd *rand.Rand) ohsnap.Arbitrary[RrrParams] {
	return Rrr(rnd, RrrAnd)
}
