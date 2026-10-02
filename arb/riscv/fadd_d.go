package riscv

// Generator for FaddD - one generator, one type, one constructor
// (over the shared floating-point three-register core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// FaddD - an arbitrary FaddD.
func FaddD(rnd *rand.Rand) ohsnap.Arbitrary[Fp3Params] {
	return Fp3(rnd, Fp3FaddD)
}
