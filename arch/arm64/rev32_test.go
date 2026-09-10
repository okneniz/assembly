package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRev32Build(t *testing.T) {
	cases := []struct {
		name string
		rd   string
		rn   string
		word uint32
	}{
		{"rev32 x0,x1", "x0", "x1", 0xdac00820},
		{"rev32 x4,xzr", "x4", "xzr", 0xdac00be4},
	}
	for _, c := range cases {
		in, err := New().Rev32(reg(t, c.rd), reg(t, c.rn))
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Rev32(reg(t, first.rd), reg(t, first.rn))
	require.NoError(t, err)
	_, ok := in.(Rev32)
	require.True(t, ok, "type = %T, want Rev32", in)

	errCases := []struct {
		name string
		rd   string
		rn   string
	}{
		{"rev32 w form", "w0", "w1"},
		{"rev32 w,rn", "x0", "w1"},
	}
	for _, c := range errCases {
		_, err := New().Rev32(reg(t, c.rd), reg(t, c.rn))
		assertErr(t, c.name, err)
	}
}
