package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSbfmBuild(t *testing.T) {
	cases := []struct {
		name string
		rd   string
		rn   string
		immr uint32
		imms uint32
		word uint32
	}{
		{"sbfm x0,x1,#5,#63 (asr)", "x0", "x1", 5, 63, 0x9345fc20},
		{"sbfm x1,x2,#0,#7 (sxtb)", "x1", "x2", 0, 7, 0x93401c41},
		{"sbfm w3,w4,#1,#3 (sbfx)", "w3", "w4", 1, 3, 0x13010c83},
	}
	for _, c := range cases {
		in, err := New().Sbfm(reg(t, c.rd), reg(t, c.rn), c.immr, c.imms)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Sbfm(reg(t, first.rd), reg(t, first.rn), first.immr, first.imms)
	require.NoError(t, err)
	_, ok := in.(Sbfm)
	require.True(t, ok, "type = %T, want Sbfm", in)

	errCases := []struct {
		name string
		rd   string
		rn   string
		immr uint32
		imms uint32
	}{
		{"sbfm sp,rd", "sp", "x1", 5, 63},
		{"sbfm sp,rn", "x0", "sp", 5, 63},
		{"sbfm x,w widths", "x0", "w1", 5, 63},
		{"sbfm immr 64", "x0", "x1", 64, 63},
		{"sbfm imms 64", "x0", "x1", 5, 64},
	}
	for _, c := range errCases {
		_, err := New().Sbfm(reg(t, c.rd), reg(t, c.rn), c.immr, c.imms)
		assertErr(t, c.name, err)
	}
}
