package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAdcBuild(t *testing.T) {
	cases := []struct {
		name string
		rd   string
		rn   string
		rm   string
		word uint32
	}{
		{"adc x0,x1,x2", "x0", "x1", "x2", 0x9a020020},
		{"adc x5,xzr,x7", "x5", "xzr", "x7", 0x9a0703e5},
	}
	for _, c := range cases {
		in, err := New().Adc(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm))
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Adc(reg(t, first.rd), reg(t, first.rn), reg(t, first.rm))
	require.NoError(t, err)
	_, ok := in.(Adc)
	require.True(t, ok, "type = %T, want Adc", in)

	errCases := []struct {
		name string
		rd   string
		rn   string
		rm   string
	}{
		{"adc w form", "w0", "x1", "x2"},
		{"adc sp", "sp", "x1", "x2"},
		{"adc rn w form", "x0", "w1", "x2"},
		{"adc rm sp", "x0", "x1", "sp"},
	}
	for _, c := range errCases {
		_, err := New().Adc(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm))
		assertErr(t, c.name, err)
	}
}
