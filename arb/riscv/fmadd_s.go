package riscv

// Generator for FmaddS - one generator, one type, one constructor
// (over the shared fused floating-point core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// FmaddS - an arbitrary FmaddS.
func FmaddS(rnd *rand.Rand) ohsnap.Arbitrary[Fp4Params] {
	return Fp4(rnd, Fp4FmaddS)
}
