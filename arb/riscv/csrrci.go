package riscv

// Generator for Csrrci - one generator, one type, one constructor
// (over the shared CSR core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// Csrrci - an arbitrary Csrrci.
func Csrrci(rnd *rand.Rand) ohsnap.Arbitrary[CsrParams] {
	return Csr(rnd, CsrCsrrci)
}
