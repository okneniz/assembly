package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEonShiftBuild(t *testing.T) {
	cases := []struct {
		name string
		rd   string
		rn   string
		rm   string
		imm  int64
		sh   Shift
		word uint32
	}{
		{"eon x1,x2,x3,ror#5", "x1", "x2", "x3", 5, ROR, 0xcae31441},
		{"eon w1,w2,w3", "w1", "w2", "w3", 0, LSL, 0x4a230041},
	}
	for _, c := range cases {
		in, err := New().EonShift(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm), imm6(t, c.imm), c.sh)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().EonShift(reg(t, first.rd), reg(t, first.rn), reg(t, first.rm), imm6(t, first.imm), first.sh)
	require.NoError(t, err)
	_, ok := in.(EonShift)
	require.True(t, ok, "type = %T, want EonShift", in)

	errCases := []struct {
		name string
		rd   string
		rn   string
		rm   string
		imm  int64
		sh   Shift
	}{
		{"eon x+w", "x0", "w1", "x2", 1, LSL},
		{"eonshift sp", "sp", "x1", "x2", 1, LSL},
		{"eonshift rn sp", "x0", "sp", "x2", 1, LSL},
		{"eonshift rm sp", "x0", "x1", "sp", 1, LSL},
	}
	for _, c := range errCases {
		_, err := New().EonShift(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm), imm6(t, c.imm), c.sh)
		assertErr(t, c.name, err)
	}
}
