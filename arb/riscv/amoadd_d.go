package riscv

// Generator for AmoaddD - one generator, one type, one constructor
// (over the shared three-register core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// AmoaddD - an arbitrary AmoaddD.
func AmoaddD(rnd *rand.Rand) ohsnap.Arbitrary[RrrParams] {
	return Rrr(rnd, RrrAmoaddD)
}
