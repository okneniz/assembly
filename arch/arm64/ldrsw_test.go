package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLdrswBuild(t *testing.T) {
	cases := []struct {
		name string
		rt   string
		rn   string
		off  Off
		word uint32
	}{
		{"ldrsw x0,[x1]", "x0", "x1", 0, 0xb9800020},
		{"ldrsw x2,[sp,#0xffc]", "x2", "sp", 0xffc, 0xb98fffe2},
		{"ldrsw xzr,[x3,#4]", "xzr", "x3", 4, 0xb980047f},
	}
	for _, c := range cases {
		in, err := New().Ldrsw(reg(t, c.rt), reg(t, c.rn), c.off)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Ldrsw(reg(t, first.rt), reg(t, first.rn), first.off)
	require.NoError(t, err)
	_, ok := in.(Ldrsw)
	require.True(t, ok, "type = %T, want Ldrsw", in)

	errCases := []struct {
		name string
		rt   string
		rn   string
		off  Off
	}{
		{"ldrsw w,rt", "w0", "x1", 0},
		{"ldrsw w1 base", "x0", "w1", 0},
		{"ldrsw off 2", "x0", "x1", 2},
		{"ldrsw off 0x10000", "x0", "x1", 0x10000},
	}
	for _, c := range errCases {
		_, err := New().Ldrsw(reg(t, c.rt), reg(t, c.rn), c.off)
		assertErr(t, c.name, err)
	}
}
