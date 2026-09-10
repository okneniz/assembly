package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAndsImmBuild(t *testing.T) {
	cases := []struct {
		name string
		rd   string
		rn   string
		imm  uint64
		word uint32
	}{
		{"tst x1,#0x7", "xzr", "x1", 0x7, 0xf240083f},
		{"ands x2,x3,#0xffff0000ffff0000", "x2", "x3", 0xffff0000ffff0000, 0xf2103c62},
		{"ands w2,w3,#0x00ff00ff", "w2", "w3", 0x00ff00ff, 0x72009c62},
	}
	for _, c := range cases {
		in, err := New().AndsImm(reg(t, c.rd), reg(t, c.rn), c.imm)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().AndsImm(reg(t, first.rd), reg(t, first.rn), first.imm)
	require.NoError(t, err)
	_, ok := in.(AndsImm)
	require.True(t, ok, "type = %T, want AndsImm", in)

	errCases := []struct {
		name string
		rd   string
		rn   string
		imm  uint64
	}{
		{"ands imm sp", "sp", "x1", 0x7},
		{"ands imm rn sp", "x0", "sp", 0x7},
		{"ands imm x+w", "x0", "w1", 0x7},
		{"ands imm not encodable", "x0", "x1", ^uint64(0)},
	}
	for _, c := range errCases {
		_, err := New().AndsImm(reg(t, c.rd), reg(t, c.rn), c.imm)
		assertErr(t, c.name, err)
	}
}
