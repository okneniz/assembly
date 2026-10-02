package riscv

// Generator for Lhu - one generator, one type, one constructor
// (over the shared load core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// Lhu - an arbitrary Lhu.
func Lhu(rnd *rand.Rand) ohsnap.Arbitrary[LoadParams] {
	return Load(rnd, LoadLhu)
}
