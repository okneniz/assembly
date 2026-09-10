package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStxrbBuild(t *testing.T) {
	cases := []struct {
		name string
		rs   string
		rt   string
		rn   string
		word uint32
	}{
		{"stxrb w0,w1,[x2]", "w0", "w1", "x2", 0x08000041},
		{"stxrb wzr,w1,[sp]", "wzr", "w1", "sp", 0x081f03e1},
	}
	for _, c := range cases {
		in, err := New().Stxrb(reg(t, c.rs), reg(t, c.rt), reg(t, c.rn))
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Stxrb(reg(t, first.rs), reg(t, first.rt), reg(t, first.rn))
	require.NoError(t, err)
	_, ok := in.(Stxrb)
	require.True(t, ok, "type = %T, want Stxrb", in)

	errCases := []struct {
		name string
		rs   string
		rt   string
		rn   string
	}{
		{"stxrb x,rs", "x0", "w1", "x2"},
		{"stxrb x,rt", "w0", "x1", "x2"},
		{"stxrb w2 base", "w0", "w1", "w2"},
	}
	for _, c := range errCases {
		_, err := New().Stxrb(reg(t, c.rs), reg(t, c.rt), reg(t, c.rn))
		assertErr(t, c.name, err)
	}
}
