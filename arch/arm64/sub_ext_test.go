package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSubExtBuild(t *testing.T) {
	cases := []struct {
		name string
		rd   string
		rn   string
		rm   string
		ext  string
		imm3 uint32
		word uint32
	}{
		{"sub x0,x1,x2,uxtx", "x0", "x1", "x2", "uxtx", 0, 0xcb226020},
		{"sub sp,sp,x2,uxtx#3", "sp", "sp", "x2", "uxtx", 3, 0xcb226fff},
		{"sub w1,w2,w3,uxtw#2", "w1", "w2", "w3", "uxtw", 2, 0x4b234841},
	}
	for _, c := range cases {
		in, err := New().SubExt(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm), c.ext, c.imm3)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().SubExt(reg(t, first.rd), reg(t, first.rn), reg(t, first.rm), first.ext, first.imm3)
	require.NoError(t, err)
	_, ok := in.(SubExt)
	require.True(t, ok, "type = %T, want SubExt", in)

	errCases := []struct {
		name string
		rd   string
		rn   string
		rm   string
		ext  string
		imm3 uint32
	}{
		{"sub xzr,rd", "xzr", "x1", "x2", "uxtx", 0},
		{"sub xzr,rn", "x0", "xzr", "x2", "uxtx", 0},
		{"sub xzr,rm", "x0", "x1", "xzr", "uxtx", 0},
		{"sub x,w widths", "x0", "w1", "x2", "uxtx", 0},
		{"sub bad ext", "x0", "x1", "x2", "foo", 0},
		{"sub imm3 8", "x0", "x1", "x2", "uxtx", 8},
	}
	for _, c := range errCases {
		_, err := New().SubExt(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm), c.ext, c.imm3)
		assertErr(t, c.name, err)
	}
}
