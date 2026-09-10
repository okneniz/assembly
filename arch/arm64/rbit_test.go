package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRbitBuild(t *testing.T) {
	cases := []struct {
		name string
		rd   string
		rn   string
		word uint32
	}{
		{"rbit x0,x1", "x0", "x1", 0xdac00020},
		{"rbit w2,w3", "w2", "w3", 0x5ac00062},
		{"rbit xzr,x30", "xzr", "x30", 0xdac003df},
	}
	for _, c := range cases {
		in, err := New().Rbit(reg(t, c.rd), reg(t, c.rn))
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Rbit(reg(t, first.rd), reg(t, first.rn))
	require.NoError(t, err)
	_, ok := in.(Rbit)
	require.True(t, ok, "type = %T, want Rbit", in)

	errCases := []struct {
		name string
		rd   string
		rn   string
	}{
		{"rbit sp,rd", "sp", "x1"},
		{"rbit sp,rn", "x0", "sp"},
		{"rbit x,w widths", "x0", "w1"},
	}
	for _, c := range errCases {
		_, err := New().Rbit(reg(t, c.rd), reg(t, c.rn))
		assertErr(t, c.name, err)
	}
}
