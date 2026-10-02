package riscv

// Generator for FnmaddD - one generator, one type, one constructor
// (over the shared fused floating-point core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// FnmaddD - an arbitrary FnmaddD.
func FnmaddD(rnd *rand.Rand) ohsnap.Arbitrary[Fp4Params] {
	return Fp4(rnd, Fp4FnmaddD)
}
