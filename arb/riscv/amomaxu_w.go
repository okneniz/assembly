package riscv

// Generator for AmomaxuW - one generator, one type, one constructor
// (over the shared three-register core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// AmomaxuW - an arbitrary AmomaxuW.
func AmomaxuW(rnd *rand.Rand) ohsnap.Arbitrary[RrrParams] {
	return Rrr(rnd, RrrAmomaxuW)
}
