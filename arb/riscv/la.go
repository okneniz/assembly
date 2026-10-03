package riscv

// Generator for the la pseudo-form - one generator, one type, one
// constructor (La over the shared target core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// La — an arbitrary la.
func La(rnd *rand.Rand) ohsnap.Arbitrary[LaParams] {
	return newLaGen(rnd)
}
