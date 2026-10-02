package riscv

// Generator for FdivS - one generator, one type, one constructor
// (over the shared floating-point three-register core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// FdivS - an arbitrary FdivS.
func FdivS(rnd *rand.Rand) ohsnap.Arbitrary[Fp3Params] {
	return Fp3(rnd, Fp3FdivS)
}
