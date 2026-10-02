package riscv

// Generator for Slliw - one generator, one type, one constructor
// (over the shared immediate-shift core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// Slliw - an arbitrary Slliw.
func Slliw(rnd *rand.Rand) ohsnap.Arbitrary[ShiftParams] {
	return Shift(rnd, ShiftSlliw)
}
