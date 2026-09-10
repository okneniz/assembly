package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAndImmBuild(t *testing.T) {
	cases := []struct {
		name string
		rd   string
		rn   string
		imm  uint64
		word uint32
	}{
		{"and x0,x1,#0x7", "x0", "x1", 0x7, 0x92400820},
		{"and x0,x1,#0xffff", "x0", "x1", 0xffff, 0x92403c20},
		{"and w0,w1,#0x00ff00ff", "w0", "w1", 0x00ff00ff, 0x12009c20},
	}
	for _, c := range cases {
		in, err := New().AndImm(reg(t, c.rd), reg(t, c.rn), c.imm)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().AndImm(reg(t, first.rd), reg(t, first.rn), first.imm)
	require.NoError(t, err)
	_, ok := in.(AndImm)
	require.True(t, ok, "type = %T, want AndImm", in)

	errCases := []struct {
		name string
		rd   string
		rn   string
		imm  uint64
	}{
		{"and imm sp", "sp", "x1", 0x7},
		{"and imm rn sp", "x0", "sp", 0x7},
		{"and imm x+w", "x0", "w1", 0x7},
		{"and imm not encodable", "x0", "x1", 0},
	}
	for _, c := range errCases {
		_, err := New().AndImm(reg(t, c.rd), reg(t, c.rn), c.imm)
		assertErr(t, c.name, err)
	}
}
