package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMaddBuild(t *testing.T) {
	cases := []struct {
		name string
		rd   string
		rn   string
		rm   string
		ra   string
		word uint32
	}{
		{"madd x1,x2,x3,x4", "x1", "x2", "x3", "x4", 0x9b031041},
		{"mul x1,x2,x3", "x1", "x2", "x3", "xzr", 0x9b037c41},
		{"madd w1,w2,w3,w4", "w1", "w2", "w3", "w4", 0x1b031041},
	}
	for _, c := range cases {
		in, err := New().Madd(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm), reg(t, c.ra))
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Madd(reg(t, first.rd), reg(t, first.rn), reg(t, first.rm), reg(t, first.ra))
	require.NoError(t, err)
	_, ok := in.(Madd)
	require.True(t, ok, "type = %T, want Madd", in)

	errCases := []struct {
		name string
		rd   string
		rn   string
		rm   string
		ra   string
	}{
		{"madd x+w", "x1", "w2", "x3", "x4"},
		{"madd sp", "x1", "x2", "x3", "sp"},
		{"madd rd sp", "sp", "x2", "x3", "x4"},
		{"madd rn sp", "x1", "sp", "x3", "x4"},
		{"madd rm sp", "x1", "x2", "sp", "x4"},
	}
	for _, c := range errCases {
		_, err := New().Madd(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm), reg(t, c.ra))
		assertErr(t, c.name, err)
	}
}
