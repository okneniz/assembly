package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSubShiftBuild(t *testing.T) {
	cases := []struct {
		name  string
		rd    string
		rn    string
		rm    string
		imm6  int64
		shift Shift
		word  uint32
	}{
		{"sub lsr#2", "x1", "x2", "x3", 2, LSR, 0xcb430841},
	}
	for _, c := range cases {
		in, err := New().SubShift(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm), imm6(t, c.imm6), c.shift)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().SubShift(reg(t, first.rd), reg(t, first.rn), reg(t, first.rm), imm6(t, first.imm6), LSL)
	require.NoError(t, err)
	_, ok := in.(SubShift)
	require.True(t, ok, "type = %T, want SubShift", in)

	errCases := []struct {
		name  string
		rd    string
		rn    string
		rm    string
		imm6  int64
		shift Shift
	}{
		{"subshift sp", "x0", "x1", "sp", 1, LSL},
		{"subshift x+w", "x0", "w1", "x2", 1, LSL},
		{"subshift w + imm6=32", "w0", "w1", "w2", 32, LSL},
	}
	for _, c := range errCases {
		_, err := New().SubShift(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm), imm6(t, c.imm6), c.shift)
		assertErr(t, c.name, err)
	}
}
