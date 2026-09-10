package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSturBuild(t *testing.T) {
	cases := []struct {
		name string
		rt   string
		rn   string
		off  Off
		word uint32
	}{
		{"stur x0,[x1]", "x0", "x1", 0, 0xf8000020},
		{"stur w1,[x2,#8]", "w1", "x2", 8, 0xb8008041},
		{"stur x3,[x4,#-256]", "x3", "x4", -256, 0xf8100083},
	}
	for _, c := range cases {
		in, err := New().Stur(reg(t, c.rt), reg(t, c.rn), c.off)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Stur(reg(t, first.rt), reg(t, first.rn), first.off)
	require.NoError(t, err)
	_, ok := in.(Stur)
	require.True(t, ok, "type = %T, want Stur", in)

	errCases := []struct {
		name string
		rt   string
		rn   string
		off  Off
	}{
		{"stur sp,rt", "sp", "x1", 0},
		{"stur w1 base", "x0", "w1", 0},
		{"stur off -257", "x0", "x1", -257},
	}
	for _, c := range errCases {
		_, err := New().Stur(reg(t, c.rt), reg(t, c.rn), c.off)
		assertErr(t, c.name, err)
	}
}
