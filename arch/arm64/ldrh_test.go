package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLdrhBuild(t *testing.T) {
	cases := []struct {
		name string
		rt   string
		rn   string
		off  Off
		word uint32
	}{
		{"ldrh w0,[x1,#2]", "w0", "x1", 2, 0x79400420},
		{"ldrh w1,[sp,#0x1ffe]", "w1", "sp", 0x1ffe, 0x797fffe1},
	}
	for _, c := range cases {
		in, err := New().Ldrh(reg(t, c.rt), reg(t, c.rn), c.off)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Ldrh(reg(t, first.rt), reg(t, first.rn), first.off)
	require.NoError(t, err)
	_, ok := in.(Ldrh)
	require.True(t, ok, "type = %T, want Ldrh", in)

	errCases := []struct {
		name string
		rt   string
		rn   string
		off  Off
	}{
		{"ldrh x,rt", "x0", "x1", 2},
		{"ldrh w1 base", "w0", "w1", 2},
		{"ldrh off 1", "w0", "x1", 1},
		{"ldrh off 0x2000", "w0", "x1", 0x2000},
	}
	for _, c := range errCases {
		_, err := New().Ldrh(reg(t, c.rt), reg(t, c.rn), c.off)
		assertErr(t, c.name, err)
	}
}
