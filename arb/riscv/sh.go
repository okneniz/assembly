package riscv

// Generator for Sh - one generator, one type, one constructor
// (over the shared store core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// Sh - an arbitrary Sh.
func Sh(rnd *rand.Rand) ohsnap.Arbitrary[StoreParams] {
	return Store(rnd, StoreSh)
}
