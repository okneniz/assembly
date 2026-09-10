package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLdrsbBuild(t *testing.T) {
	cases := []struct {
		name string
		rt   string
		rn   string
		off  Off
		word uint32
	}{
		{"ldrsb x0,[x1]", "x0", "x1", 0, 0x39800020},
		{"ldrsb x2,[sp,#0xfff]", "x2", "sp", 0xfff, 0x39bfffe2},
	}
	for _, c := range cases {
		in, err := New().Ldrsb(reg(t, c.rt), reg(t, c.rn), c.off)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Ldrsb(reg(t, first.rt), reg(t, first.rn), first.off)
	require.NoError(t, err)
	_, ok := in.(Ldrsb)
	require.True(t, ok, "type = %T, want Ldrsb", in)

	errCases := []struct {
		name string
		rt   string
		rn   string
		off  Off
	}{
		{"ldrsb w,rt", "w0", "x1", 0},
		{"ldrsb w1 base", "x0", "w1", 0},
		{"ldrsb off -1", "x0", "x1", -1},
	}
	for _, c := range errCases {
		_, err := New().Ldrsb(reg(t, c.rt), reg(t, c.rn), c.off)
		assertErr(t, c.name, err)
	}
}
