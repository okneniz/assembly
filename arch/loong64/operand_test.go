package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestR(t *testing.T) {
	for n := range 32 {
		r, err := R(n)
		require.NoError(t, err, "R(%d)", n)
		require.Equal(t, uint8(n), r.Num())
	}

	for _, n := range []int{-1, 32, 255} {
		_, err := R(n)
		require.ErrorContains(t, err, "outside 0..31", "R(%d)", n)
	}
}

func TestRegNames(t *testing.T) {
	for _, tc := range []struct {
		num  int
		name string
	}{
		{0, "$zero"},
		{1, "$ra"},
		{2, "$tp"},
		{3, "$sp"},
		{4, "$a0"},
		{11, "$a7"},
		{12, "$t0"},
		{20, "$t8"},
		{21, "$r21"},
		{22, "$fp"},
		{23, "$s0"},
		{31, "$s8"},
	} {
		require.Equal(t, tc.name, lreg(t, tc.num).String(), "R(%d)", tc.num)
	}

	require.Equal(t, "$zero", Zero.String())
	require.Equal(t, "$ra", Ra.String())
	require.Equal(t, "$tp", Tp.String())
	require.Equal(t, "$sp", Sp.String())
	require.Equal(t, "$fp", Fp.String())
}
