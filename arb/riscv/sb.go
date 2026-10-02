package riscv

// Generator for Sb - one generator, one type, one constructor
// (over the shared store core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// Sb - an arbitrary Sb.
func Sb(rnd *rand.Rand) ohsnap.Arbitrary[StoreParams] {
	return Store(rnd, StoreSb)
}
