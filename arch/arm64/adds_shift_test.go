package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddsShiftBuild(t *testing.T) {
	cases := []struct {
		name string
		rd   string
		rn   string
		rm   string
		imm  int64
		sh   Shift
		word uint32
	}{
		{"cmn x1,x2", "xzr", "x1", "x2", 0, LSL, 0xab02003f},
		{"adds x1,x2,x3,lsl#4", "x1", "x2", "x3", 4, LSL, 0xab031041},
		{"adds w1,w2,w3,asr#5", "w1", "w2", "w3", 5, ASR, 0x2b831441},
		{"adds x1,x2,x3,lsr#63", "x1", "x2", "x3", 63, LSR, 0xab43fc41},
	}
	for _, c := range cases {
		in, err := New().AddsShift(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm), imm6(t, c.imm), c.sh)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().AddsShift(reg(t, first.rd), reg(t, first.rn), reg(t, first.rm), imm6(t, first.imm), first.sh)
	require.NoError(t, err)
	_, ok := in.(AddsShift)
	require.True(t, ok, "type = %T, want AddsShift", in)

	errCases := []struct {
		name string
		rd   string
		rn   string
		rm   string
		imm  int64
		sh   Shift
	}{
		{"adds x+w", "x0", "w1", "x2", 1, LSL},
		{"addsshift sp", "sp", "x1", "x2", 1, LSL},
		{"addsshift rn sp", "x0", "sp", "x2", 1, LSL},
		{"addsshift rm sp", "x0", "x1", "sp", 1, LSL},
		{"adds ror", "x0", "x1", "x2", 1, ROR},
		{"adds w + imm6=32", "w0", "w1", "w2", 32, LSL},
	}
	for _, c := range errCases {
		_, err := New().AddsShift(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm), imm6(t, c.imm), c.sh)
		assertErr(t, c.name, err)
	}
}
