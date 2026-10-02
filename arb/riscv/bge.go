package riscv

// Generator for Bge - one generator, one type, one constructor
// (over the shared branch core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// Bge - an arbitrary Bge.
func Bge(rnd *rand.Rand) ohsnap.Arbitrary[BranchParams] {
	return Branch(rnd, BranchBge)
}
