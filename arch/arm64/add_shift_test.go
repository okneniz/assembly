package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddShiftBuild(t *testing.T) {
	cases := []struct {
		name string
		rd   string
		rn   string
		rm   string
		imm  int64
		sh   Shift
		word uint32
	}{
		{"add x1,x2,x3", "x1", "x2", "x3", 0, LSL, 0x8b030041},
		{"add x1,x2,x3,lsl#4", "x1", "x2", "x3", 4, LSL, 0x8b031041},
		{"add w1,w2,w3,asr#5", "w1", "w2", "w3", 5, ASR, 0x0b831441},
	}
	for _, c := range cases {
		in, err := New().AddShift(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm), imm6(t, c.imm), c.sh)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().AddShift(reg(t, first.rd), reg(t, first.rn), reg(t, first.rm), imm6(t, first.imm), first.sh)
	require.NoError(t, err)
	_, ok := in.(AddShift)
	require.True(t, ok, "type = %T, want AddShift", in)

	errCases := []struct {
		name string
		rd   string
		rn   string
		rm   string
		imm  int64
		sh   Shift
	}{
		{"add x+w", "x0", "w1", "x2", 1, LSL},
		{"addshift sp", "sp", "x1", "x2", 1, LSL},
		{"add ror", "x0", "x1", "x2", 1, ROR},
		{"add w + imm6=32", "w0", "w1", "w2", 32, LSL},
	}
	for _, c := range errCases {
		_, err := New().AddShift(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm), imm6(t, c.imm), c.sh)
		assertErr(t, c.name, err)
	}
}
