package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUbfmBuild(t *testing.T) {
	cases := []struct {
		name string
		rd   string
		rn   string
		immr uint32
		imms uint32
		word uint32
	}{
		{"ubfm x0,x1,#5,#63 (lsr)", "x0", "x1", 5, 63, 0xd345fc20},
		{"ubfm x2,x3,#59,#58 (lsl)", "x2", "x3", 59, 58, 0xd37be862},
		{"ubfm w1,w2,#0,#31 (lsr w)", "w1", "w2", 0, 31, 0x53007c41},
	}
	for _, c := range cases {
		in, err := New().Ubfm(reg(t, c.rd), reg(t, c.rn), c.immr, c.imms)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Ubfm(reg(t, first.rd), reg(t, first.rn), first.immr, first.imms)
	require.NoError(t, err)
	_, ok := in.(Ubfm)
	require.True(t, ok, "type = %T, want Ubfm", in)

	errCases := []struct {
		name string
		rd   string
		rn   string
		immr uint32
		imms uint32
	}{
		{"ubfm sp,rd", "sp", "x1", 5, 63},
		{"ubfm sp,rn", "x0", "sp", 5, 63},
		{"ubfm x,w widths", "x0", "w1", 5, 63},
		{"ubfm immr 64", "x0", "x1", 64, 63},
		{"ubfm imms 64", "x0", "x1", 5, 64},
	}
	for _, c := range errCases {
		_, err := New().Ubfm(reg(t, c.rd), reg(t, c.rn), c.immr, c.imms)
		assertErr(t, c.name, err)
	}
}
