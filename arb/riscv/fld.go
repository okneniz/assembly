package riscv

// Generator for Fld - one generator, one type, one constructor
// (over the shared load core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// Fld - an arbitrary Fld.
func Fld(rnd *rand.Rand) ohsnap.Arbitrary[LoadParams] {
	return Load(rnd, LoadFld)
}
