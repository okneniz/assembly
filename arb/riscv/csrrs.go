package riscv

// Generator for Csrrs - one generator, one type, one constructor
// (over the shared CSR core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// Csrrs - an arbitrary Csrrs.
func Csrrs(rnd *rand.Rand) ohsnap.Arbitrary[CsrParams] {
	return Csr(rnd, CsrCsrrs)
}
