package riscv

// Generator for Andi - one generator, one type, one constructor
// (over the shared I-type core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// Andi - an arbitrary Andi.
func Andi(rnd *rand.Rand) ohsnap.Arbitrary[RiParams] {
	return Ri(rnd, RiAndi)
}
