package riscv

// Generator for FnmaddS - one generator, one type, one constructor
// (over the shared fused floating-point core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// FnmaddS - an arbitrary FnmaddS.
func FnmaddS(rnd *rand.Rand) ohsnap.Arbitrary[Fp4Params] {
	return Fp4(rnd, Fp4FnmaddS)
}
