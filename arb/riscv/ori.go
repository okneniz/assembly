package riscv

// Generator for Ori - one generator, one type, one constructor
// (over the shared I-type core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// Ori - an arbitrary Ori.
func Ori(rnd *rand.Rand) ohsnap.Arbitrary[RiParams] {
	return Ri(rnd, RiOri)
}
