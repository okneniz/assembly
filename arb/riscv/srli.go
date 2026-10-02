package riscv

// Generator for Srli - one generator, one type, one constructor
// (over the shared immediate-shift core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// Srli - an arbitrary Srli.
func Srli(rnd *rand.Rand) ohsnap.Arbitrary[ShiftParams] {
	return Shift(rnd, ShiftSrli)
}
