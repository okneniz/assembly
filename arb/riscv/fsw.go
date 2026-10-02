package riscv

// Generator for Fsw - one generator, one type, one constructor
// (over the shared store core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// Fsw - an arbitrary Fsw.
func Fsw(rnd *rand.Rand) ohsnap.Arbitrary[StoreParams] {
	return Store(rnd, StoreFsw)
}
