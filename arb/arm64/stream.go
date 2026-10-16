package arm64

// The local alias of the package-wide helper (arb.Enum):
// one definition point for the generators of this package.

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
)

// arbEnum — the Enum adapter of this package.
func arbEnum[T comparable](rnd *rand.Rand, vals []T) ohsnap.Arbitrary[T] {
	return arb.Enum(rnd, vals...)
}
