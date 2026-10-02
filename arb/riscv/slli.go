package riscv

// Generator for Slli - one generator, one type, one constructor
// (over the shared immediate-shift core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// Slli - an arbitrary Slli.
func Slli(rnd *rand.Rand) ohsnap.Arbitrary[ShiftParams] {
	return Shift(rnd, ShiftSlli)
}
