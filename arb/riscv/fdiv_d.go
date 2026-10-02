package riscv

// Generator for FdivD - one generator, one type, one constructor
// (over the shared floating-point three-register core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// FdivD - an arbitrary FdivD.
func FdivD(rnd *rand.Rand) ohsnap.Arbitrary[Fp3Params] {
	return Fp3(rnd, Fp3FdivD)
}
