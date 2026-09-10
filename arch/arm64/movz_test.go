package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMovzBuild(t *testing.T) {
	cases := []struct {
		name string
		rd   string
		imm  int64
		hw   Hw
		word uint32
	}{
		{"movz x0,#0x1234", "x0", 0x1234, Hw0, 0xd2824680},
		{"movz x0,#1", "x0", 1, Hw0, 0xd2800020},
		{"movz w3,#0x42,hw1", "w3", 0x42, Hw1, 0x52a00843},
		{"movz xzr,#0,hw3", "xzr", 0, Hw3, 0xd2e0001f},
	}
	for _, c := range cases {
		in, err := New().Movz(reg(t, c.rd), imm16(t, c.imm), c.hw)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Movz(reg(t, first.rd), imm16(t, first.imm), first.hw)
	require.NoError(t, err)
	_, ok := in.(Movz)
	require.True(t, ok, "type = %T, want Movz", in)

	errCases := []struct {
		name string
		rd   string
		imm  int64
		hw   Hw
	}{
		{"movz sp", "sp", 1, Hw0},
		{"movz w0,hw2", "w0", 1, Hw2},
	}
	for _, c := range errCases {
		_, err := New().Movz(reg(t, c.rd), imm16(t, c.imm), c.hw)
		assertErr(t, c.name, err)
	}
}
