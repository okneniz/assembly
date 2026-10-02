package riscv

// Generator for Blt - one generator, one type, one constructor
// (over the shared branch core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// Blt - an arbitrary Blt.
func Blt(rnd *rand.Rand) ohsnap.Arbitrary[BranchParams] {
	return Branch(rnd, BranchBlt)
}
