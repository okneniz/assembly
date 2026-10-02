package riscv

// Generator for AmoorD - one generator, one type, one constructor
// (over the shared three-register core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// AmoorD - an arbitrary AmoorD.
func AmoorD(rnd *rand.Rand) ohsnap.Arbitrary[RrrParams] {
	return Rrr(rnd, RrrAmoorD)
}
