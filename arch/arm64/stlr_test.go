package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStlrBuild(t *testing.T) {
	cases := []struct {
		name string
		rt   string
		rn   string
		word uint32
	}{
		{"stlr x0,[x1]", "x0", "x1", 0xc89ffc20},
		{"stlr w2,[sp]", "w2", "sp", 0x889fffe2},
		{"stlr xzr,[x3]", "xzr", "x3", 0xc89ffc7f},
	}
	for _, c := range cases {
		in, err := New().Stlr(reg(t, c.rt), reg(t, c.rn))
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Stlr(reg(t, first.rt), reg(t, first.rn))
	require.NoError(t, err)
	_, ok := in.(Stlr)
	require.True(t, ok, "type = %T, want Stlr", in)

	errCases := []struct {
		name string
		rt   string
		rn   string
	}{
		{"stlr sp,rt", "sp", "x1"},
		{"stlr w1 base", "x0", "w1"},
	}
	for _, c := range errCases {
		_, err := New().Stlr(reg(t, c.rt), reg(t, c.rn))
		assertErr(t, c.name, err)
	}
}
