package riscv

// Generator for FnmsubD - one generator, one type, one constructor
// (over the shared fused floating-point core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// FnmsubD - an arbitrary FnmsubD.
func FnmsubD(rnd *rand.Rand) ohsnap.Arbitrary[Fp4Params] {
	return Fp4(rnd, Fp4FnmsubD)
}
