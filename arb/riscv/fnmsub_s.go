package riscv

// Generator for FnmsubS - one generator, one type, one constructor
// (over the shared fused floating-point core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// FnmsubS - an arbitrary FnmsubS.
func FnmsubS(rnd *rand.Rand) ohsnap.Arbitrary[Fp4Params] {
	return Fp4(rnd, Fp4FnmsubS)
}
