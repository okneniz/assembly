package riscv

// Generator for Fsd - one generator, one type, one constructor
// (over the shared store core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// Fsd - an arbitrary Fsd.
func Fsd(rnd *rand.Rand) ohsnap.Arbitrary[StoreParams] {
	return Store(rnd, StoreFsd)
}
