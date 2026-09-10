package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddsExtBuild(t *testing.T) {
	cases := []struct {
		name string
		rd   string
		rn   string
		rm   string
		ext  string
		imm3 uint32
		word uint32
	}{
		{"cmn x1,x2,uxtb#0", "xzr", "x1", "x2", "uxtb", 0, 0xab22003f},
		{"adds w1,w2,w3,sxtw#1", "w1", "w2", "w3", "sxtw", 1, 0x2b23c441},
		{"adds x1,sp,x3,uxtx#0", "x1", "sp", "x3", "uxtx", 0, 0xab2363e1},
	}
	for _, c := range cases {
		in, err := New().AddsExt(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm), c.ext, c.imm3)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().AddsExt(reg(t, first.rd), reg(t, first.rn), reg(t, first.rm), first.ext, first.imm3)
	require.NoError(t, err)
	_, ok := in.(AddsExt)
	require.True(t, ok, "type = %T, want AddsExt", in)

	errCases := []struct {
		name string
		rd   string
		rn   string
		rm   string
		ext  string
		imm3 uint32
	}{
		{"adds ext rd sp", "sp", "x1", "x2", "uxtb", 0},
		{"adds ext rn xzr", "x0", "xzr", "x2", "uxtb", 0},
		{"adds ext rm xzr", "x0", "x1", "xzr", "uxtb", 0},
		{"adds ext x+w", "x0", "w1", "x2", "uxtb", 0},
		{"adds ext bad ext", "x0", "x1", "x2", "sxtb2", 0},
		{"adds ext imm3=8", "x0", "x1", "x2", "uxtx", 8},
	}
	for _, c := range errCases {
		_, err := New().AddsExt(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm), c.ext, c.imm3)
		assertErr(t, c.name, err)
	}
}
