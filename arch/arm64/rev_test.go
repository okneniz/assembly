package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRevBuild(t *testing.T) {
	cases := []struct {
		name string
		rd   string
		rn   string
		word uint32
	}{
		{"rev x0,x1", "x0", "x1", 0xdac00c20},
		{"rev xzr,x2", "xzr", "x2", 0xdac00c5f},
	}
	for _, c := range cases {
		in, err := New().Rev(reg(t, c.rd), reg(t, c.rn))
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Rev(reg(t, first.rd), reg(t, first.rn))
	require.NoError(t, err)
	_, ok := in.(Rev)
	require.True(t, ok, "type = %T, want Rev", in)

	errCases := []struct {
		name string
		rd   string
		rn   string
	}{
		{"rev w form", "w0", "w1"},
		{"rev w,rn", "x0", "w1"},
	}
	for _, c := range errCases {
		_, err := New().Rev(reg(t, c.rd), reg(t, c.rn))
		assertErr(t, c.name, err)
	}
}
