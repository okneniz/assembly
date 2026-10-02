package riscv

// Generator for Beq - one generator, one type, one constructor
// (over the shared branch core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// Beq - an arbitrary Beq.
func Beq(rnd *rand.Rand) ohsnap.Arbitrary[BranchParams] {
	return Branch(rnd, BranchBeq)
}
