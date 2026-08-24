// Package arb — arbitrary value generators (oh-snap Arbitrary) for
// assembly property tests. Architecture-specific generators live in
// the subpackages arb/arm64 and arb/riscv (one generator per instruction,
// on top of the arch/* constructors); here are the common utilities.
//
// Seed: determinism is a mandatory condition for reproducible failures;
// Rnd(seed) gives a PCG generator, tests read the seed from ASSEMBLY_SEED.
package arb

import (
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

// Halved — shrink candidates for a number by halving toward zero (v/2, v/4, ..., 0).
// Sign-preserving: negative values approach zero through negatives.
func Halved(v int64) []int64 {
	var out []int64
	for d := v / 2; d != 0; d /= 2 {
		out = append(out, d)
	}

	if v != 0 {
		out = append(out, 0)
	}

	return out
}
