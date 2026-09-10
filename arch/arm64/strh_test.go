package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStrhBuild(t *testing.T) {
	cases := []struct {
		name string
		rt   string
		rn   string
		off  Off
		word uint32
	}{
		{"strh w0,[x1]", "w0", "x1", 0, 0x79000020},
		{"strh w2,[x3,#0x1ffe]", "w2", "x3", 0x1ffe, 0x793ffc62},
	}
	for _, c := range cases {
		in, err := New().Strh(reg(t, c.rt), reg(t, c.rn), c.off)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Strh(reg(t, first.rt), reg(t, first.rn), first.off)
	require.NoError(t, err)
	_, ok := in.(Strh)
	require.True(t, ok, "type = %T, want Strh", in)

	errCases := []struct {
		name string
		rt   string
		rn   string
		off  Off
	}{
		{"strh x,rt", "x0", "x1", 0},
		{"strh w1 base", "w0", "w1", 0},
		{"strh off 1", "w0", "x1", 1},
	}
	for _, c := range errCases {
		_, err := New().Strh(reg(t, c.rt), reg(t, c.rn), c.off)
		assertErr(t, c.name, err)
	}
}
