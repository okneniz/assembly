package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEorImmBuild(t *testing.T) {
	cases := []struct {
		name string
		rd   string
		rn   string
		imm  uint64
		word uint32
	}{
		{"eor x0,x1,#0x7", "x0", "x1", 0x7, 0xd2400820},
		{"eor x2,x3,#0xffff", "x2", "x3", 0xffff, 0xd2403c62},
		{"eor w0,w1,#0x00ff00ff", "w0", "w1", 0x00ff00ff, 0x52009c20},
	}
	for _, c := range cases {
		in, err := New().EorImm(reg(t, c.rd), reg(t, c.rn), c.imm)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().EorImm(reg(t, first.rd), reg(t, first.rn), first.imm)
	require.NoError(t, err)
	_, ok := in.(EorImm)
	require.True(t, ok, "type = %T, want EorImm", in)

	errCases := []struct {
		name string
		rd   string
		rn   string
		imm  uint64
	}{
		{"eor imm sp", "sp", "x1", 0x7},
		{"eor imm rn sp", "x0", "sp", 0x7},
		{"eor imm x+w", "x0", "w1", 0x7},
		{"eor imm not encodable", "x0", "x1", 0x55},
	}
	for _, c := range errCases {
		_, err := New().EorImm(reg(t, c.rd), reg(t, c.rn), c.imm)
		assertErr(t, c.name, err)
	}
}
