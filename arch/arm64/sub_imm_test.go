package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSubImmBuild(t *testing.T) {
	cases := []struct {
		name  string
		rd    string
		rn    string
		imm12 int64
		sh    Sh12
		word  uint32
	}{
		{"sub", "x0", "x1", 0x42, NoSh12, 0xd1010820},
		{"sub w0,wsp", "w0", "wsp", 1, NoSh12, 0x510007e0},
	}
	for _, c := range cases {
		in, err := New().SubImm(reg(t, c.rd), reg(t, c.rn), imm12(t, c.imm12), c.sh)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().SubImm(reg(t, first.rd), reg(t, first.rn), imm12(t, first.imm12), NoSh12)
	require.NoError(t, err)
	_, ok := in.(SubImm)
	require.True(t, ok, "type = %T, want SubImm", in)

	errCases := []struct {
		name  string
		rd    string
		rn    string
		imm12 int64
	}{
		{"sub imm zr", "xzr", "x1", 1},
		{"sub imm rn zr", "x0", "xzr", 1},
		{"sub imm x+w", "x0", "w1", 1},
	}
	for _, c := range errCases {
		_, err := New().SubImm(reg(t, c.rd), reg(t, c.rn), imm12(t, c.imm12), NoSh12)
		assertErr(t, c.name, err)
	}
}
