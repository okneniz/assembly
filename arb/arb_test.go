package arb

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestWordShrink — shrinking a word by halving toward zero (built-in ArbitraryUint32).
func TestWordShrink(t *testing.T) {
	w := Word(Rnd(42))
	require.NotEmpty(t, w.Shrink(w.Generate()), "Word.Shrink is empty")
}

// TestHalved — sign-preserving shrink helper.
func TestHalved(t *testing.T) {
	got := Halved(9)
	require.Len(t, got, 4, "Halved(9) = %v", got)
	require.Equal(t, int64(4), got[0], "Halved(9) = %v", got)
	require.Equal(t, int64(0), got[3], "Halved(9) = %v", got)
	got = Halved(-9)
	require.Len(t, got, 4, "Halved(-9) = %v", got)
	require.Equal(t, int64(-4), got[0], "Halved(-9) = %v", got)
	require.Equal(t, int64(0), got[3], "Halved(-9) = %v", got)
	require.Empty(t, Halved(0), "Halved(0)")
	got = Halved(1)
	require.Len(t, got, 1, "Halved(1) = %v", got)
	require.Equal(t, int64(0), got[0], "Halved(1) = %v", got)
}

// TestRndDeterminism — the same seed yields the same sequence.
func TestRndDeterminism(t *testing.T) {
	a, b := Rnd(7), Rnd(7)
	for range 100 {
		require.Equal(t, a.Uint32(), b.Uint32(), "sequences diverged")
	}
}
