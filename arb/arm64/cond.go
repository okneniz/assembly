package arm64

// The condition generator: uniform over the arch table, shrink to "eq".

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/arm64"
)

// Cond — an arbitrary condition name (the full arch table; the alias
// families that forbid al/nv filter on their side).
func Cond(rnd *rand.Rand) ohsnap.Arbitrary[string] {
	names := arm64.CondNames()
	vals := make([]string, 0, len(names))
	for _, name := range names {
		vals = append(vals, name)
	}

	return arb.Enum(rnd, vals...)
}
