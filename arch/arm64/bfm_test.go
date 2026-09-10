package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBfmBuild(t *testing.T) {
	cases := []struct {
		name string
		rd   string
		rn   string
		immr uint32
		imms uint32
		word uint32
	}{
		{"bfm x0,x1,#5,#7", "x0", "x1", 5, 7, 0xb3451c20},
		{"bfm w2,w3,#0,#31", "w2", "w3", 0, 31, 0x33007c62},
		{"bfm x0,x1,#63,#63", "x0", "x1", 63, 63, 0xb37ffc20},
	}
	for _, c := range cases {
		in, err := New().Bfm(reg(t, c.rd), reg(t, c.rn), c.immr, c.imms)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Bfm(reg(t, first.rd), reg(t, first.rn), first.immr, first.imms)
	require.NoError(t, err)
	_, ok := in.(Bfm)
	require.True(t, ok, "type = %T, want Bfm", in)

	errCases := []struct {
		name string
		rd   string
		rn   string
		immr uint32
		imms uint32
	}{
		{"bfm sp,rd", "sp", "x1", 5, 7},
		{"bfm sp,rn", "x0", "sp", 5, 7},
		{"bfm x,w widths", "x0", "w1", 5, 7},
		{"bfm immr 64", "x0", "x1", 64, 7},
		{"bfm imms 64", "x0", "x1", 5, 64},
	}
	for _, c := range errCases {
		_, err := New().Bfm(reg(t, c.rd), reg(t, c.rn), c.immr, c.imms)
		assertErr(t, c.name, err)
	}
}
