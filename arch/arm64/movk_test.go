package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMovkBuild(t *testing.T) {
	cases := []struct {
		name  string
		rd    string
		imm16 int64
		hw    Hw
		word  uint32
	}{
		{"movk", "x0", 0x1234, Hw1, 0xf2a24680},
		{"movk w4", "w4", 0x1234, Hw0, 0x72824684},
	}
	for _, c := range cases {
		in, err := New().Movk(reg(t, c.rd), imm16(t, c.imm16), c.hw)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Movk(reg(t, first.rd), imm16(t, first.imm16), Hw0)
	require.NoError(t, err)
	_, ok := in.(Movk)
	require.True(t, ok, "type = %T, want Movk", in)

	errCases := []struct {
		name  string
		rd    string
		imm16 int64
		hw    Hw
	}{
		{"movk w0,hw3", "w0", 1, Hw3},
		{"movk wsp", "wsp", 1, Hw0},
	}
	for _, c := range errCases {
		_, err := New().Movk(reg(t, c.rd), imm16(t, c.imm16), c.hw)
		assertErr(t, c.name, err)
	}
}
