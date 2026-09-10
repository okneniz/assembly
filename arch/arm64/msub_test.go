package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMsubBuild(t *testing.T) {
	cases := []struct {
		name string
		rd   string
		rn   string
		rm   string
		ra   string
		word uint32
	}{
		{"msub x1,x2,x3,x4", "x1", "x2", "x3", "x4", 0x9b039041},
		{"mneg w1,w2,w3", "w1", "w2", "w3", "wzr", 0x1b03fc41},
	}
	for _, c := range cases {
		in, err := New().Msub(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm), reg(t, c.ra))
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Msub(reg(t, first.rd), reg(t, first.rn), reg(t, first.rm), reg(t, first.ra))
	require.NoError(t, err)
	_, ok := in.(Msub)
	require.True(t, ok, "type = %T, want Msub", in)

	errCases := []struct {
		name string
		rd   string
		rn   string
		rm   string
		ra   string
	}{
		{"msub x+w", "x1", "w2", "x3", "x4"},
		{"msub sp", "x1", "x2", "x3", "sp"},
	}
	for _, c := range errCases {
		_, err := New().Msub(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm), reg(t, c.ra))
		assertErr(t, c.name, err)
	}
}
