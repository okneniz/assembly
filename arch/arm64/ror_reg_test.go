package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRorRegBuild(t *testing.T) {
	cases := []struct {
		name string
		rd   string
		rn   string
		rm   string
		word uint32
	}{
		{"ror x0,x1,x2", "x0", "x1", "x2", 0x9a022c20},
		{"ror x3,xzr,x4", "x3", "xzr", "x4", 0x9a042fe3},
	}
	for _, c := range cases {
		in, err := New().RorReg(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm))
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().RorReg(reg(t, first.rd), reg(t, first.rn), reg(t, first.rm))
	require.NoError(t, err)
	_, ok := in.(RorReg)
	require.True(t, ok, "type = %T, want RorReg", in)

	errCases := []struct {
		name string
		rd   string
		rn   string
		rm   string
	}{
		{"ror w form", "w0", "w1", "w2"},
		{"ror w,rn", "x0", "w1", "x2"},
		{"ror w,rm", "x0", "x1", "w2"},
	}
	for _, c := range errCases {
		_, err := New().RorReg(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm))
		assertErr(t, c.name, err)
	}
}
