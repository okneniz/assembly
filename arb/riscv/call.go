package riscv

// Generator for the call pseudo-form - one generator, one type, one
// constructor (Call over the shared target core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// Call — an arbitrary call.
func Call(rnd *rand.Rand) ohsnap.Arbitrary[CallParams] {
	return newCallGen(rnd)
}
