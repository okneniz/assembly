package riscv

// Generator for Srliw - one generator, one type, one constructor
// (over the shared immediate-shift core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// Srliw - an arbitrary Srliw.
func Srliw(rnd *rand.Rand) ohsnap.Arbitrary[ShiftParams] {
	return Shift(rnd, ShiftSrliw)
}
