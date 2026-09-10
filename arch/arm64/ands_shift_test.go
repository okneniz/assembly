package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAndsShiftBuild(t *testing.T) {
	cases := []struct {
		name string
		rd   string
		rn   string
		rm   string
		imm  int64
		sh   Shift
		word uint32
	}{
		{"tst w1,w2", "wzr", "w1", "w2", 0, LSL, 0x6a02003f},
		{"ands x1,x2,x3,lsl#4", "x1", "x2", "x3", 4, LSL, 0xea031041},
		{"ands w1,w2,w3,asr#5", "w1", "w2", "w3", 5, ASR, 0x6a831441},
	}
	for _, c := range cases {
		in, err := New().AndsShift(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm), imm6(t, c.imm), c.sh)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().AndsShift(reg(t, first.rd), reg(t, first.rn), reg(t, first.rm), imm6(t, first.imm), first.sh)
	require.NoError(t, err)
	_, ok := in.(AndsShift)
	require.True(t, ok, "type = %T, want AndsShift", in)

	errCases := []struct {
		name string
		rd   string
		rn   string
		rm   string
		imm  int64
		sh   Shift
	}{
		{"ands x+w", "x0", "w1", "x2", 1, LSL},
		{"andsshift sp", "sp", "x1", "x2", 1, LSL},
		{"andsshift rn sp", "x0", "sp", "x2", 1, LSL},
		{"andsshift rm sp", "x0", "x1", "sp", 1, LSL},
	}
	for _, c := range errCases {
		_, err := New().AndsShift(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm), imm6(t, c.imm), c.sh)
		assertErr(t, c.name, err)
	}
}
