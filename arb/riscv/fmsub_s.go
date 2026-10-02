package riscv

// Generator for FmsubS - one generator, one type, one constructor
// (over the shared fused floating-point core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// FmsubS - an arbitrary FmsubS.
func FmsubS(rnd *rand.Rand) ohsnap.Arbitrary[Fp4Params] {
	return Fp4(rnd, Fp4FmsubS)
}
