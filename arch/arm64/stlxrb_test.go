package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStlxrbBuild(t *testing.T) {
	cases := []struct {
		name string
		rs   string
		rt   string
		rn   string
		word uint32
	}{
		{"stlxrb w0,w1,[x2]", "w0", "w1", "x2", 0x0800fc41},
		{"stlxrb wzr,w0,[x1]", "wzr", "w0", "x1", 0x081ffc20},
	}
	for _, c := range cases {
		in, err := New().Stlxrb(reg(t, c.rs), reg(t, c.rt), reg(t, c.rn))
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Stlxrb(reg(t, first.rs), reg(t, first.rt), reg(t, first.rn))
	require.NoError(t, err)
	_, ok := in.(Stlxrb)
	require.True(t, ok, "type = %T, want Stlxrb", in)

	errCases := []struct {
		name string
		rs   string
		rt   string
		rn   string
	}{
		{"stlxrb x,rs", "x0", "w1", "x2"},
		{"stlxrb x,rt", "w0", "x1", "x2"},
		{"stlxrb w2 base", "w0", "w1", "w2"},
	}
	for _, c := range errCases {
		_, err := New().Stlxrb(reg(t, c.rs), reg(t, c.rt), reg(t, c.rn))
		assertErr(t, c.name, err)
	}
}
