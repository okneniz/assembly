package riscv

// Generator for FsubS - one generator, one type, one constructor
// (over the shared floating-point three-register core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// FsubS - an arbitrary FsubS.
func FsubS(rnd *rand.Rand) ohsnap.Arbitrary[Fp3Params] {
	return Fp3(rnd, Fp3FsubS)
}
