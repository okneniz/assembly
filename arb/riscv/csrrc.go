package riscv

// Generator for Csrrc - one generator, one type, one constructor
// (over the shared CSR core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// Csrrc - an arbitrary Csrrc.
func Csrrc(rnd *rand.Rand) ohsnap.Arbitrary[CsrParams] {
	return Csr(rnd, CsrCsrrc)
}
