package riscv

// Generator for Slti - one generator, one type, one constructor
// (over the shared I-type core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// Slti - an arbitrary Slti.
func Slti(rnd *rand.Rand) ohsnap.Arbitrary[RiParams] {
	return Ri(rnd, RiSlti)
}
