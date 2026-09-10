package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddExtBuild(t *testing.T) {
	cases := []struct {
		name string
		rd   string
		rn   string
		rm   string
		ext  string
		imm3 uint32
		word uint32
	}{
		{"add x1,x2,x3,uxtb#0", "x1", "x2", "x3", "uxtb", 0, 0x8b230041},
		{"add x0,x1,x2,uxtx#7", "x0", "x1", "x2", "uxtx", 7, 0x8b227c20},
		{"add w1,w2,w3,uxtw#2", "w1", "w2", "w3", "uxtw", 2, 0x0b234841},
		{"add sp,sp,x3,sxtx#3", "sp", "sp", "x3", "sxtx", 3, 0x8b23efff},
	}
	for _, c := range cases {
		in, err := New().AddExt(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm), c.ext, c.imm3)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().AddExt(reg(t, first.rd), reg(t, first.rn), reg(t, first.rm), first.ext, first.imm3)
	require.NoError(t, err)
	_, ok := in.(AddExt)
	require.True(t, ok, "type = %T, want AddExt", in)

	errCases := []struct {
		name string
		rd   string
		rn   string
		rm   string
		ext  string
		imm3 uint32
	}{
		{"add ext rd xzr", "xzr", "x1", "x2", "uxtb", 0},
		{"add ext rn xzr", "x0", "xzr", "x2", "uxtb", 0},
		{"add ext rm xzr", "x0", "x1", "xzr", "uxtb", 0},
		{"add ext x+w", "x0", "w1", "x2", "uxtb", 0},
		{"add ext bad ext", "x0", "x1", "x2", "uxtx2", 0},
		{"add ext imm3=8", "x0", "x1", "x2", "uxtx", 8},
	}
	for _, c := range errCases {
		_, err := New().AddExt(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm), c.ext, c.imm3)
		assertErr(t, c.name, err)
	}
}
