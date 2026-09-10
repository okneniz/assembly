package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBicsShiftBuild(t *testing.T) {
	cases := []struct {
		name string
		rd   string
		rn   string
		rm   string
		imm  int64
		sh   Shift
		word uint32
	}{
		{"bics x1,x2,x3", "x1", "x2", "x3", 0, LSL, 0xea230041},
		{"bics w1,w2,w3,lsl#4", "w1", "w2", "w3", 4, LSL, 0x6a231041},
		{"bics x1,xzr,x3,ror#2", "x1", "xzr", "x3", 2, ROR, 0xeae30be1},
	}
	for _, c := range cases {
		in, err := New().BicsShift(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm), imm6(t, c.imm), c.sh)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().BicsShift(reg(t, first.rd), reg(t, first.rn), reg(t, first.rm), imm6(t, first.imm), first.sh)
	require.NoError(t, err)
	_, ok := in.(BicsShift)
	require.True(t, ok, "type = %T, want BicsShift", in)

	errCases := []struct {
		name string
		rd   string
		rn   string
		rm   string
		imm  int64
		sh   Shift
	}{
		{"bics x+w", "x0", "w1", "x2", 1, LSL},
		{"bicsshift sp", "sp", "x1", "x2", 1, LSL},
		{"bicsshift rn sp", "x0", "sp", "x2", 1, LSL},
		{"bicsshift rm sp", "x0", "x1", "sp", 1, LSL},
	}
	for _, c := range errCases {
		_, err := New().BicsShift(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm), imm6(t, c.imm), c.sh)
		assertErr(t, c.name, err)
	}
}
