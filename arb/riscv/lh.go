package riscv

// Generator for Lh - one generator, one type, one constructor
// (over the shared load core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// Lh - an arbitrary Lh.
func Lh(rnd *rand.Rand) ohsnap.Arbitrary[LoadParams] {
	return Load(rnd, LoadLh)
}
