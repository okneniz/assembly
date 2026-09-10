package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTbzBuild(t *testing.T) {
	cases := []struct {
		name   string
		rt     string
		bit    uint32
		target int64
		word   uint32
	}{
		{"tbz w0,#5,0x1000", "w0", 5, 0, 0x36280000},
		{"tbz w1,#0,0x1010", "w1", 0, 16, 0x36000081},
		{"tbz w0,#31,0x8ffc", "w0", 31, 32764, 0x36fbffe0},
		{"tbz x2,#32,0x1000", "x2", 32, 0, 0xb6000002},
		{"tbz w0,#1,0xf00", "w0", 1, -256, 0x360ff800},
	}
	for _, c := range cases {
		in, err := New().Tbz(reg(t, c.rt), c.bit, c.target)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Tbz(reg(t, first.rt), first.bit, first.target)
	require.NoError(t, err)
	_, ok := in.(Tbz)
	require.True(t, ok, "type = %T, want Tbz", in)

	errCases := []struct {
		name   string
		rt     string
		bit    uint32
		target int64
	}{
		{"tbz sp,rt", "sp", 5, 0x1000},
		{"tbz bit 64", "x0", 64, 0x1000},
		{"tbz x0,#5", "x0", 5, 0x1000},
		{"tbz w0,#32", "w0", 32, 0x1000},
	}
	for _, c := range errCases {
		_, err := New().Tbz(reg(t, c.rt), c.bit, c.target)
		assertErr(t, c.name, err)
	}
}
