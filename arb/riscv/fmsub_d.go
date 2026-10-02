package riscv

// Generator for FmsubD - one generator, one type, one constructor
// (over the shared fused floating-point core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// FmsubD - an arbitrary FmsubD.
func FmsubD(rnd *rand.Rand) ohsnap.Arbitrary[Fp4Params] {
	return Fp4(rnd, Fp4FmsubD)
}
