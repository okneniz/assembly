package riscv

// Generator for FmulS - one generator, one type, one constructor
// (over the shared floating-point three-register core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// FmulS - an arbitrary FmulS.
func FmulS(rnd *rand.Rand) ohsnap.Arbitrary[Fp3Params] {
	return Fp3(rnd, Fp3FmulS)
}
