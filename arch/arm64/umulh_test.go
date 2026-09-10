package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUmulhBuild(t *testing.T) {
	cases := []struct {
		name string
		rd   string
		rn   string
		rm   string
		word uint32
	}{
		{"umulh x3,x4,x5", "x3", "x4", "x5", 0x9bc57c83},
		{"umulh x0,xzr,x2", "x0", "xzr", "x2", 0x9bc27fe0},
	}
	for _, c := range cases {
		in, err := New().Umulh(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm))
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Umulh(reg(t, first.rd), reg(t, first.rn), reg(t, first.rm))
	require.NoError(t, err)
	_, ok := in.(Umulh)
	require.True(t, ok, "type = %T, want Umulh", in)

	errCases := []struct {
		name string
		rd   string
		rn   string
		rm   string
	}{
		{"umulh w form", "w3", "w4", "w5"},
		{"umulh w,rn", "x3", "w4", "x5"},
		{"umulh w,rm", "x3", "x4", "w5"},
	}
	for _, c := range errCases {
		_, err := New().Umulh(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm))
		assertErr(t, c.name, err)
	}
}
