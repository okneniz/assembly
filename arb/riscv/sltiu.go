package riscv

// Generator for Sltiu - one generator, one type, one constructor
// (over the shared I-type core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// Sltiu - an arbitrary Sltiu.
func Sltiu(rnd *rand.Rand) ohsnap.Arbitrary[RiParams] {
	return Ri(rnd, RiSltiu)
}
