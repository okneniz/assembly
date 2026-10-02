package riscv

// Generator for Csrrw - one generator, one type, one constructor
// (over the shared CSR core).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// Csrrw - an arbitrary Csrrw.
func Csrrw(rnd *rand.Rand) ohsnap.Arbitrary[CsrParams] {
	return Csr(rnd, CsrCsrrw)
}
