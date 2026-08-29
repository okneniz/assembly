// Package arb — arbitrary value generators (oh-snap Arbitrary) for
// assembly property tests. Architecture-specific generators live in
// the subpackages arb/arm64 and arb/riscv (one generator per instruction,
// on top of the arch/* constructors); here are the common utilities.
//
// Seed: determinism is a mandatory condition for reproducible failures;
// Rnd(seed) gives a PCG generator, tests read the seed from ASSEMBLY_SEED.
package arb

import (
	"iter"
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

// Rnd — a deterministic generator (math/rand/v2, PCG) from a seed.
func Rnd(seed uint64) *rand.Rand {
	return rand.New(rand.NewPCG(0, seed))
}

// Word — an arbitrary 32-bit word (robustness properties of decoders).
func Word(rnd *rand.Rand) ohsnap.Arbitrary[uint32] {
	return ohsnap.ArbitraryUint32(rnd, 0, 0xffffffff)
}

// Stream — a lazy infinite sequence of values from f: the adapter that
// turns a plain value constructor into the oh-snap Generate sequence.
func Stream[T any](f func() T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for {
			if !yield(f()) {
				return
			}
		}
	}
}
