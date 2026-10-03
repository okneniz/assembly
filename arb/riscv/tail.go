package riscv

// Generator for the tail pseudo-form - one generator, one type, one
// constructor (Tail over the shared target core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// Tail — an arbitrary tail.
func Tail(rnd *rand.Rand) ohsnap.Arbitrary[TailParams] {
	return newTailGen(rnd)
}
