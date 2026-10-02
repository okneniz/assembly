package riscv

// Generator for Lb - one generator, one type, one constructor
// (over the shared load core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// Lb - an arbitrary Lb.
func Lb(rnd *rand.Rand) ohsnap.Arbitrary[LoadParams] {
	return Load(rnd, LoadLb)
}
