package arb

import (
	"slices"
	"testing"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/stretchr/testify/require"
)

// TestWordShrink — shrinking a word by halving toward zero (built-in ArbitraryUint32).
func TestWordShrink(t *testing.T) {
	w := Word(Rnd(42))
	candidates := slices.Collect(w.Shrink(ohsnap.First(w.Generate())))
	require.NotEmpty(t, candidates, "Word.Shrink is empty")
}

// TestStream — the sequence yields f() values until the consumer stops:
// after an early break f is not called anymore (laziness).
func TestStream(t *testing.T) {
	calls := 0
	seq := Stream(func() int {
		calls++
		return calls
	})

	got := slices.Collect(func(yield func(int) bool) {
		for v := range seq {
			if v == 2 {
				return
			}

			if !yield(v) {
				return
			}
		}
	})

	require.Equal(t, []int{1}, got, "Stream yielded %v", got)
	require.Equal(t, 2, calls, "Stream called f after the consumer stopped")
}

// TestRndDeterminism — the same seed yields the same sequence.
func TestRndDeterminism(t *testing.T) {
	a, b := Rnd(7), Rnd(7)
	for range 100 {
		require.Equal(t, a.Uint32(), b.Uint32(), "sequences diverged")
	}
}
