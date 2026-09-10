package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBlrBuild(t *testing.T) {
	cases := []struct {
		name string
		rn   string
		word uint32
	}{
		{"blr x1", "x1", 0xd63f0020},
		{"blr x30", "x30", 0xd63f03c0},
		{"blr xzr", "xzr", 0xd63f03e0},
	}
	for _, c := range cases {
		in, err := New().Blr(reg(t, c.rn))
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Blr(reg(t, first.rn))
	require.NoError(t, err)
	_, ok := in.(Blr)
	require.True(t, ok, "type = %T, want Blr", in)

	errCases := []struct {
		name string
		rn   string
	}{
		{"blr w1", "w1"},
	}
	for _, c := range errCases {
		_, err := New().Blr(reg(t, c.rn))
		assertErr(t, c.name, err)
	}
}
