package riscv

// Generator for Bne - one generator, one type, one constructor
// (over the shared branch core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// Bne - an arbitrary Bne.
func Bne(rnd *rand.Rand) ohsnap.Arbitrary[BranchParams] {
	return Branch(rnd, BranchBne)
}
