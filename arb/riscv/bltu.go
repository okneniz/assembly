package riscv

// Generator for Bltu - one generator, one type, one constructor
// (over the shared branch core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// Bltu - an arbitrary Bltu.
func Bltu(rnd *rand.Rand) ohsnap.Arbitrary[BranchParams] {
	return Branch(rnd, BranchBltu)
}
