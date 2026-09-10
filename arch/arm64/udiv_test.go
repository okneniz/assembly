package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUdivBuild(t *testing.T) {
	cases := []struct {
		name string
		rd   string
		rn   string
		rm   string
		word uint32
	}{
		{"udiv x0,x1,x2", "x0", "x1", "x2", 0x9ac20820},
		{"udiv w4,w5,w6", "w4", "w5", "w6", 0x1ac608a4},
	}
	for _, c := range cases {
		in, err := New().Udiv(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm))
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Udiv(reg(t, first.rd), reg(t, first.rn), reg(t, first.rm))
	require.NoError(t, err)
	_, ok := in.(Udiv)
	require.True(t, ok, "type = %T, want Udiv", in)

	errCases := []struct {
		name string
		rd   string
		rn   string
		rm   string
	}{
		{"udiv sp,rd", "sp", "x1", "x2"},
		{"udiv sp,rn", "x0", "sp", "x2"},
		{"udiv sp,rm", "x0", "x1", "sp"},
		{"udiv x,w widths", "x0", "w1", "x2"},
	}
	for _, c := range errCases {
		_, err := New().Udiv(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm))
		assertErr(t, c.name, err)
	}
}
