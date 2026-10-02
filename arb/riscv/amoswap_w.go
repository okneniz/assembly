package riscv

// Generator for AmoswapW - one generator, one type, one constructor
// (over the shared three-register core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// AmoswapW - an arbitrary AmoswapW.
func AmoswapW(rnd *rand.Rand) ohsnap.Arbitrary[RrrParams] {
	return Rrr(rnd, RrrAmoswapW)
}
