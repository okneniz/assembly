package arb

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"
)

// enumArb — enum generator: uniform over the values, shrink to the first.
type enumArb[T comparable] struct {
	rnd  *rand.Rand
	vals []T
}

// Enum — an arbitrary value from the listed ones. Shrink collapses to vals[0]:
// in a minimal counterexample the enum axis is canonicalized (Hw→Hw0, Shift→LSL);
// a value that "failed to shrink" without losing the failure is itself a suspect.
func Enum[T comparable](rnd *rand.Rand, vals ...T) ohsnap.Arbitrary[T] {
	return enumArb[T]{
		rnd:  rnd,
		vals: vals,
	}
}

func (a enumArb[T]) Generate() iter.Seq[T] {
	return Stream(func() T {
		return a.vals[a.rnd.IntN(len(a.vals))]
	})
}

func (a enumArb[T]) Shrink(v T) iter.Seq[T] {
	if v == a.vals[0] {
		return ohsnap.Empty[T]()
	}

	return slices.Values([]T{a.vals[0]})
}
