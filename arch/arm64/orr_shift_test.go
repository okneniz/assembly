package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOrrShiftBuild(t *testing.T) {
	cases := []struct {
		name string
		rd   string
		rn   string
		rm   string
		imm  int64
		sh   Shift
		word uint32
	}{
		{"orr x1,x2,x3,lsl#4", "x1", "x2", "x3", 4, LSL, 0xaa031041},
		{"orr w1,w2,w3", "w1", "w2", "w3", 0, LSL, 0x2a030041},
		{"mov x1,x3,ror#2", "x1", "xzr", "x3", 2, ROR, 0xaac30be1},
	}
	for _, c := range cases {
		in, err := New().OrrShift(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm), imm6(t, c.imm), c.sh)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().OrrShift(reg(t, first.rd), reg(t, first.rn), reg(t, first.rm), imm6(t, first.imm), first.sh)
	require.NoError(t, err)
	_, ok := in.(OrrShift)
	require.True(t, ok, "type = %T, want OrrShift", in)

	errCases := []struct {
		name string
		rd   string
		rn   string
		rm   string
		imm  int64
		sh   Shift
	}{
		{"orr x+w", "x0", "w1", "x2", 1, LSL},
		{"orrshift sp", "x0", "x1", "sp", 1, LSL},
		{"orrshift rd sp", "sp", "x1", "x2", 1, LSL},
		{"orrshift rn sp", "x0", "sp", "x2", 1, LSL},
	}
	for _, c := range errCases {
		_, err := New().OrrShift(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm), imm6(t, c.imm), c.sh)
		assertErr(t, c.name, err)
	}
}
