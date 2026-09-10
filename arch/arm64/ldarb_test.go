package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLdarbBuild(t *testing.T) {
	cases := []struct {
		name string
		rt   string
		rn   string
		word uint32
	}{
		{"ldarb w3,[x4]", "w3", "x4", 0x08dffc83},
		{"ldarb wzr,[x29]", "wzr", "x29", 0x08dfffbf},
	}
	for _, c := range cases {
		in, err := New().Ldarb(reg(t, c.rt), reg(t, c.rn))
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Ldarb(reg(t, first.rt), reg(t, first.rn))
	require.NoError(t, err)
	_, ok := in.(Ldarb)
	require.True(t, ok, "type = %T, want Ldarb", in)

	errCases := []struct {
		name string
		rt   string
		rn   string
	}{
		{"ldarb x,rt", "x3", "x4"},
		{"ldarb w1 base", "w3", "w4"},
	}
	for _, c := range errCases {
		_, err := New().Ldarb(reg(t, c.rt), reg(t, c.rn))
		assertErr(t, c.name, err)
	}
}
