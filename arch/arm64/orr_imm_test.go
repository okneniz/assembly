package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOrrImmBuild(t *testing.T) {
	cases := []struct {
		name string
		rd   string
		rn   string
		imm  uint64
		word uint32
	}{
		{"orr x0,x1,#0x7", "x0", "x1", 0x7, 0xb2400820},
		{"orr x0,xzr,#0xffff", "x0", "xzr", 0xffff, 0xb2403fe0},
		{"orr w0,w1,#0x00ff00ff", "w0", "w1", 0x00ff00ff, 0x32009c20},
	}
	for _, c := range cases {
		in, err := New().OrrImm(reg(t, c.rd), reg(t, c.rn), c.imm)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().OrrImm(reg(t, first.rd), reg(t, first.rn), first.imm)
	require.NoError(t, err)
	_, ok := in.(OrrImm)
	require.True(t, ok, "type = %T, want OrrImm", in)

	errCases := []struct {
		name string
		rd   string
		rn   string
		imm  uint64
	}{
		{"orr imm sp", "sp", "x1", 0x7},
		{"orr imm rn sp", "x0", "sp", 0x7},
		{"orr imm x+w", "x0", "w1", 0x7},
		{"orr imm not encodable", "x0", "x1", 0},
	}
	for _, c := range errCases {
		_, err := New().OrrImm(reg(t, c.rd), reg(t, c.rn), c.imm)
		assertErr(t, c.name, err)
	}
}
