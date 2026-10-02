package riscv

// Generator for Srai - one generator, one type, one constructor
// (over the shared immediate-shift core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// Srai - an arbitrary Srai.
func Srai(rnd *rand.Rand) ohsnap.Arbitrary[ShiftParams] {
	return Shift(rnd, ShiftSrai)
}
