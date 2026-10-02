package riscv

// Generator for Xori - one generator, one type, one constructor
// (over the shared I-type core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// Xori - an arbitrary Xori.
func Xori(rnd *rand.Rand) ohsnap.Arbitrary[RiParams] {
	return Ri(rnd, RiXori)
}
