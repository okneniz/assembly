package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStpBuild(t *testing.T) {
	cases := []struct {
		name string
		rt   string
		rt2  string
		rn   string
		off  Off
		word uint32
	}{
		{"stp x0,x1,[x2]", "x0", "x1", "x2", 0, 0xa9000440},
		{"stp x29,x30,[sp,#-16]", "x29", "x30", "sp", -16, 0xa93f7bfd},
		{"stp w0,w1,[x2,#252]", "w0", "w1", "x2", 252, 0x291f8440},
	}
	for _, c := range cases {
		in, err := New().Stp(reg(t, c.rt), reg(t, c.rt2), reg(t, c.rn), c.off)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Stp(reg(t, first.rt), reg(t, first.rt2), reg(t, first.rn), first.off)
	require.NoError(t, err)
	_, ok := in.(Stp)
	require.True(t, ok, "type = %T, want Stp", in)

	errCases := []struct {
		name string
		rt   string
		rt2  string
		rn   string
		off  Off
	}{
		{"stp sp,rt", "sp", "x1", "x2", 0},
		{"stp sp,rt2", "x0", "sp", "x2", 0},
		{"stp w2 base", "x0", "x1", "w2", 0},
		{"stp x,w widths", "x0", "w1", "x2", 0},
		{"stp off 4 in x form", "x0", "x1", "x2", 4},
		{"stp off -520 in x form", "x0", "x1", "x2", -520},
	}
	for _, c := range errCases {
		_, err := New().Stp(reg(t, c.rt), reg(t, c.rt2), reg(t, c.rn), c.off)
		assertErr(t, c.name, err)
	}
}
