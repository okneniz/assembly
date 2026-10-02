package riscv

// Generator for AmomaxW - one generator, one type, one constructor
// (over the shared three-register core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// AmomaxW - an arbitrary AmomaxW.
func AmomaxW(rnd *rand.Rand) ohsnap.Arbitrary[RrrParams] {
	return Rrr(rnd, RrrAmomaxW)
}
