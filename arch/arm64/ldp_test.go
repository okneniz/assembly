package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLdpBuild(t *testing.T) {
	cases := []struct {
		name string
		rt   string
		rt2  string
		rn   string
		off  Off
		word uint32
	}{
		{"ldp x0,x1,[x2]", "x0", "x1", "x2", 0, 0xa9400440},
		{"ldp x29,x30,[sp]", "x29", "x30", "sp", 0, 0xa9407bfd},
		{"ldp w3,w4,[x5,#8]", "w3", "w4", "x5", 8, 0x294110a3},
		{"ldp x0,x1,[x2,#-512]", "x0", "x1", "x2", -512, 0xa9600440},
	}
	for _, c := range cases {
		in, err := New().Ldp(reg(t, c.rt), reg(t, c.rt2), reg(t, c.rn), c.off)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Ldp(reg(t, first.rt), reg(t, first.rt2), reg(t, first.rn), first.off)
	require.NoError(t, err)
	_, ok := in.(Ldp)
	require.True(t, ok, "type = %T, want Ldp", in)

	errCases := []struct {
		name string
		rt   string
		rt2  string
		rn   string
		off  Off
	}{
		{"ldp sp,rt", "sp", "x1", "x2", 0},
		{"ldp sp,rt2", "x0", "sp", "x2", 0},
		{"ldp w2 base", "x0", "x1", "w2", 0},
		{"ldp x,w widths", "x0", "w1", "x2", 0},
		{"ldp off 4 in x form", "x0", "x1", "x2", 4},
		{"ldp off 520 in x form", "x0", "x1", "x2", 520},
	}
	for _, c := range errCases {
		_, err := New().Ldp(reg(t, c.rt), reg(t, c.rt2), reg(t, c.rn), c.off)
		assertErr(t, c.name, err)
	}
}
