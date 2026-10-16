package arb

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestWordShrink — shrinking a word by halving toward zero (built-in arbitrary.Uint32).
func TestWordShrink(t *testing.T) {
	w := Word(Rnd(42))
	candidates := slices.Collect(w.Shrink(w.Generate()))
	require.NotEmpty(t, candidates, "Word.Shrink is empty")
}

// TestRndDeterminism — the same seed yields the same sequence.
func TestRndDeterminism(t *testing.T) {
	a, b := Rnd(7), Rnd(7)
	for range 100 {
		require.Equal(t, a.Uint32(), b.Uint32(), "sequences diverged")
	}
}
