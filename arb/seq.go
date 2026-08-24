package arb

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// seqArb — sequence generator: length 1..32, a random generator from the
// set for each position; shrink is dropping.
type seqArb[T any] struct {
	rnd  *rand.Rand
	gens []func() T
}

// Seq — an arbitrary sequence of values (1..32) from a set of
// generators. Shrink drops: half of the prefix and "drop one" candidates —
// a minimal failing subsequence localizes the bug.
// (Element-wise shrink is unavailable here: the generators are of different
// types, see Map in ohsnap.)
func Seq[T any](rnd *rand.Rand, gens []func() T) ohsnap.Arbitrary[[]T] {
	return seqArb[T]{
		rnd:  rnd,
		gens: gens,
	}
}

func (s seqArb[T]) Generate() []T {
	n := 1 + s.rnd.IntN(32)
	out := make([]T, n)
	for i := range out {
		out[i] = s.gens[s.rnd.IntN(len(s.gens))]()
	}

	return out
}

func (s seqArb[T]) Shrink(v []T) [][]T {
	var out [][]T
	if h := len(v) / 2; h > 0 {
		out = append(out, v[:h])
	}

	for i := range v {
		if len(v) == 1 {
			break
		}

		out = append(out, append(append([]T{}, v[:i]...), v[i+1:]...))
	}

	return out
}
