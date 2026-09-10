package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMovnBuild(t *testing.T) {
	cases := []struct {
		name string
		rd   string
		imm  int64
		hw   Hw
		word uint32
	}{
		{"movn x0,#0", "x0", 0, Hw0, 0x92800000},
		{"movn x1,#0x1234,Hw1", "x1", 0x1234, Hw1, 0x92a24681},
		{"movn w2,#0xffff", "w2", 0xffff, Hw0, 0x129fffe2},
		{"movn x0,#0xffff,Hw3", "x0", 0xffff, Hw3, 0x92ffffe0},
	}
	for _, c := range cases {
		in, err := New().Movn(reg(t, c.rd), imm16(t, c.imm), c.hw)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Movn(reg(t, first.rd), imm16(t, first.imm), first.hw)
	require.NoError(t, err)
	_, ok := in.(Movn)
	require.True(t, ok, "type = %T, want Movn", in)

	errCases := []struct {
		name string
		rd   string
		imm  int64
		hw   Hw
	}{
		{"movn sp,rd", "sp", 0, Hw0},
		{"movn w,Hw2", "w0", 0, Hw2},
	}
	for _, c := range errCases {
		_, err := New().Movn(reg(t, c.rd), imm16(t, c.imm), c.hw)
		assertErr(t, c.name, err)
	}
}
