package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRetBuild(t *testing.T) {
	cases := []struct {
		name string
		rn   string
		word uint32
	}{
		{"ret", "x30", 0xd65f03c0},
		{"ret x8", "x8", 0xd65f0100},
	}
	for _, c := range cases {
		in, err := New().Ret(reg(t, c.rn))
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Ret(reg(t, first.rn))
	require.NoError(t, err)
	_, ok := in.(Ret)
	require.True(t, ok, "type = %T, want Ret", in)

	errCases := []struct {
		name string
		rn   string
	}{
		{"ret sp", "sp"},
		{"ret w0", "w0"},
	}
	for _, c := range errCases {
		_, err := New().Ret(reg(t, c.rn))
		assertErr(t, c.name, err)
	}
}
