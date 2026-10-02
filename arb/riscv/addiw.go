package riscv

// Generator for Addiw - one generator, one type, one constructor
// (over the shared I-type core); the imm=0 shape is the sext.w
// pseudo-form of the decoder canon.

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// Addiw - an arbitrary Addiw.
func Addiw(rnd *rand.Rand) ohsnap.Arbitrary[RiParams] {
	return Ri(rnd, RiAddiw)
}
