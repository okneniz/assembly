package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAndShiftBuild(t *testing.T) {
	cases := []struct {
		name string
		rd   string
		rn   string
		rm   string
		imm  int64
		sh   Shift
		word uint32
	}{
		{"and x1,x2,x3", "x1", "x2", "x3", 0, LSL, 0x8a030041},
		{"and x1,x2,x3,lsl#63", "x1", "x2", "x3", 63, LSL, 0x8a03fc41},
		{"and w1,w2,w3,ror#5", "w1", "w2", "w3", 5, ROR, 0x0ac31441},
	}
	for _, c := range cases {
		in, err := New().AndShift(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm), imm6(t, c.imm), c.sh)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().AndShift(reg(t, first.rd), reg(t, first.rn), reg(t, first.rm), imm6(t, first.imm), first.sh)
	require.NoError(t, err)
	_, ok := in.(AndShift)
	require.True(t, ok, "type = %T, want AndShift", in)

	errCases := []struct {
		name string
		rd   string
		rn   string
		rm   string
		imm  int64
		sh   Shift
	}{
		{"and x+w", "x0", "w1", "x2", 1, LSL},
		{"andshift sp", "sp", "x1", "x2", 1, LSL},
		{"andshift rn sp", "x0", "sp", "x2", 1, LSL},
		{"andshift rm sp", "x0", "x1", "sp", 1, LSL},
	}
	for _, c := range errCases {
		_, err := New().AndShift(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm), imm6(t, c.imm), c.sh)
		assertErr(t, c.name, err)
	}
}
