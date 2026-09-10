package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBicShiftBuild(t *testing.T) {
	cases := []struct {
		name string
		rd   string
		rn   string
		rm   string
		imm  int64
		sh   Shift
		word uint32
	}{
		{"bic x1,x2,x3", "x1", "x2", "x3", 0, LSL, 0x8a230041},
		{"bic w1,w2,w3,asr#5", "w1", "w2", "w3", 5, ASR, 0x0aa31441},
		{"bic x1,x2,x3,lsr#63", "x1", "x2", "x3", 63, LSR, 0x8a63fc41},
	}
	for _, c := range cases {
		in, err := New().BicShift(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm), imm6(t, c.imm), c.sh)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().BicShift(reg(t, first.rd), reg(t, first.rn), reg(t, first.rm), imm6(t, first.imm), first.sh)
	require.NoError(t, err)
	_, ok := in.(BicShift)
	require.True(t, ok, "type = %T, want BicShift", in)

	errCases := []struct {
		name string
		rd   string
		rn   string
		rm   string
		imm  int64
		sh   Shift
	}{
		{"bic x+w", "x0", "w1", "x2", 1, LSL},
		{"bicshift sp", "sp", "x1", "x2", 1, LSL},
		{"bicshift rn sp", "x0", "sp", "x2", 1, LSL},
		{"bicshift rm sp", "x0", "x1", "sp", 1, LSL},
	}
	for _, c := range errCases {
		_, err := New().BicShift(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm), imm6(t, c.imm), c.sh)
		assertErr(t, c.name, err)
	}
}
