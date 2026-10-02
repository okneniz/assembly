package riscv

// Generator for AmomaxD - one generator, one type, one constructor
// (over the shared three-register core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// AmomaxD - an arbitrary AmomaxD.
func AmomaxD(rnd *rand.Rand) ohsnap.Arbitrary[RrrParams] {
	return Rrr(rnd, RrrAmomaxD)
}
