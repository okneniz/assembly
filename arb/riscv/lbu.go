package riscv

// Generator for Lbu - one generator, one type, one constructor
// (over the shared load core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// Lbu - an arbitrary Lbu.
func Lbu(rnd *rand.Rand) ohsnap.Arbitrary[LoadParams] {
	return Load(rnd, LoadLbu)
}
