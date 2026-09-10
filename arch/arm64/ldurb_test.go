package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLdurbBuild(t *testing.T) {
	cases := []struct {
		name string
		rt   string
		rn   string
		off  Off
		word uint32
	}{
		{"ldurb w0,[x1]", "w0", "x1", 0, 0x38400020},
		{"ldurb w3,[x4,#-256]", "w3", "x4", -256, 0x38500083},
	}
	for _, c := range cases {
		in, err := New().Ldurb(reg(t, c.rt), reg(t, c.rn), c.off)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Ldurb(reg(t, first.rt), reg(t, first.rn), first.off)
	require.NoError(t, err)
	_, ok := in.(Ldurb)
	require.True(t, ok, "type = %T, want Ldurb", in)

	errCases := []struct {
		name string
		rt   string
		rn   string
		off  Off
	}{
		{"ldurb x,rt", "x0", "x1", 0},
		{"ldurb w1 base", "w0", "w1", 0},
		{"ldurb off 256", "w0", "x1", 256},
	}
	for _, c := range errCases {
		_, err := New().Ldurb(reg(t, c.rt), reg(t, c.rn), c.off)
		assertErr(t, c.name, err)
	}
}
