package riscv

// Generator for AmomaxuD - one generator, one type, one constructor
// (over the shared three-register core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// AmomaxuD - an arbitrary AmomaxuD.
func AmomaxuD(rnd *rand.Rand) ohsnap.Arbitrary[RrrParams] {
	return Rrr(rnd, RrrAmomaxuD)
}
