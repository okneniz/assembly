package riscv

// Generator for AmoxorW - one generator, one type, one constructor
// (over the shared three-register core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// AmoxorW - an arbitrary AmoxorW.
func AmoxorW(rnd *rand.Rand) ohsnap.Arbitrary[RrrParams] {
	return Rrr(rnd, RrrAmoxorW)
}
