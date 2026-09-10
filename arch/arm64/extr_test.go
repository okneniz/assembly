package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtrBuild(t *testing.T) {
	cases := []struct {
		name string
		rd   string
		rn   string
		rm   string
		lsb  int64
		word uint32
	}{
		{"extr x1,x2,x3,#5", "x1", "x2", "x3", 5, 0x93031441},
		{"ror x1,x2,#4", "x1", "x2", "x2", 4, 0x93021041},
		{"extr x1,x2,x3,#63", "x1", "x2", "x3", 63, 0x9303fc41},
	}
	for _, c := range cases {
		in, err := New().Extr(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm), imm6(t, c.lsb))
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Extr(reg(t, first.rd), reg(t, first.rn), reg(t, first.rm), imm6(t, first.lsb))
	require.NoError(t, err)
	_, ok := in.(Extr)
	require.True(t, ok, "type = %T, want Extr", in)

	errCases := []struct {
		name string
		rd   string
		rn   string
		rm   string
		lsb  int64
	}{
		{"extr w form", "w1", "x2", "x3", 5},
		{"extr sp", "x1", "sp", "x3", 5},
		{"extr rm sp", "x1", "x2", "sp", 5},
	}
	for _, c := range errCases {
		_, err := New().Extr(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm), imm6(t, c.lsb))
		assertErr(t, c.name, err)
	}
}
