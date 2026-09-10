package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSubsExtBuild(t *testing.T) {
	cases := []struct {
		name string
		rd   string
		rn   string
		rm   string
		ext  string
		imm3 uint32
		word uint32
	}{
		{"subs x0,x1,x2,sxtx#2", "x0", "x1", "x2", "sxtx", 2, 0xeb22e820},
		{"subs xzr,x1,x2,sxtx (cmp)", "xzr", "x1", "x2", "sxtx", 0, 0xeb22e03f},
		{"subs w1,w2,w3,uxth", "w1", "w2", "w3", "uxth", 0, 0x6b232041},
	}
	for _, c := range cases {
		in, err := New().SubsExt(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm), c.ext, c.imm3)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().SubsExt(reg(t, first.rd), reg(t, first.rn), reg(t, first.rm), first.ext, first.imm3)
	require.NoError(t, err)
	_, ok := in.(SubsExt)
	require.True(t, ok, "type = %T, want SubsExt", in)

	errCases := []struct {
		name string
		rd   string
		rn   string
		rm   string
		ext  string
		imm3 uint32
	}{
		{"subs sp,rd", "sp", "x1", "x2", "sxtx", 2},
		{"subs xzr,rn", "x0", "xzr", "x2", "sxtx", 2},
		{"subs xzr,rm", "x0", "x1", "xzr", "sxtx", 2},
		{"subs x,w widths", "x0", "w1", "x2", "sxtx", 2},
		{"subs bad ext", "x0", "x1", "x2", "foo", 2},
		{"subs imm3 8", "x0", "x1", "x2", "sxtx", 8},
	}
	for _, c := range errCases {
		_, err := New().SubsExt(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm), c.ext, c.imm3)
		assertErr(t, c.name, err)
	}
}
