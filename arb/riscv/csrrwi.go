package riscv

// Generator for Csrrwi - one generator, one type, one constructor
// (over the shared CSR core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// Csrrwi - an arbitrary Csrrwi.
func Csrrwi(rnd *rand.Rand) ohsnap.Arbitrary[CsrParams] {
	return Csr(rnd, CsrCsrrwi)
}
