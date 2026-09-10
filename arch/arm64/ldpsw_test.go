package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLdpswBuild(t *testing.T) {
	cases := []struct {
		name string
		rt   string
		rt2  string
		rn   string
		off  Off
		word uint32
	}{
		{"ldpsw x0,x1,[x2,#4]", "x0", "x1", "x2", 4, 0x69408440},
		{"ldpsw x2,x3,[x4,#-256]", "x2", "x3", "x4", -256, 0x69600c82},
	}
	for _, c := range cases {
		in, err := New().Ldpsw(reg(t, c.rt), reg(t, c.rt2), reg(t, c.rn), c.off)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Ldpsw(reg(t, first.rt), reg(t, first.rt2), reg(t, first.rn), first.off)
	require.NoError(t, err)
	_, ok := in.(Ldpsw)
	require.True(t, ok, "type = %T, want Ldpsw", in)

	errCases := []struct {
		name string
		rt   string
		rt2  string
		rn   string
		off  Off
	}{
		{"ldpsw w,rt", "w0", "x1", "x2", 4},
		{"ldpsw w,rt2", "x0", "w1", "x2", 4},
		{"ldpsw w2 base", "x0", "x1", "w2", 4},
		{"ldpsw off 2", "x0", "x1", "x2", 2},
		{"ldpsw off 260", "x0", "x1", "x2", 260},
	}
	for _, c := range errCases {
		_, err := New().Ldpsw(reg(t, c.rt), reg(t, c.rt2), reg(t, c.rn), c.off)
		assertErr(t, c.name, err)
	}
}
