package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBrBuild(t *testing.T) {
	cases := []struct {
		name string
		rn   string
		word uint32
	}{
		{"br x5", "x5", 0xd61f00a0},
		{"br x0", "x0", 0xd61f0000},
		{"br xzr", "xzr", 0xd61f03e0},
	}
	for _, c := range cases {
		in, err := New().Br(reg(t, c.rn))
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Br(reg(t, first.rn))
	require.NoError(t, err)
	_, ok := in.(Br)
	require.True(t, ok, "type = %T, want Br", in)

	errCases := []struct {
		name string
		rn   string
	}{
		{"br w5", "w5"},
	}
	for _, c := range errCases {
		_, err := New().Br(reg(t, c.rn))
		assertErr(t, c.name, err)
	}
}
