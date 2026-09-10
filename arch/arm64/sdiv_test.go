package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSdivBuild(t *testing.T) {
	cases := []struct {
		name string
		rd   string
		rn   string
		rm   string
		word uint32
	}{
		{"sdiv x0,x1,x2", "x0", "x1", "x2", 0x9ac20c20},
		{"sdiv w4,w5,w6", "w4", "w5", "w6", 0x1ac60ca4},
	}
	for _, c := range cases {
		in, err := New().Sdiv(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm))
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Sdiv(reg(t, first.rd), reg(t, first.rn), reg(t, first.rm))
	require.NoError(t, err)
	_, ok := in.(Sdiv)
	require.True(t, ok, "type = %T, want Sdiv", in)

	errCases := []struct {
		name string
		rd   string
		rn   string
		rm   string
	}{
		{"sdiv sp,rd", "sp", "x1", "x2"},
		{"sdiv sp,rn", "x0", "sp", "x2"},
		{"sdiv sp,rm", "x0", "x1", "sp"},
		{"sdiv x,w widths", "x0", "w1", "x2"},
	}
	for _, c := range errCases {
		_, err := New().Sdiv(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm))
		assertErr(t, c.name, err)
	}
}
