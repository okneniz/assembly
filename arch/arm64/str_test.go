package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStrBuild(t *testing.T) {
	cases := []struct {
		name string
		rt   string
		rn   string
		off  Off
		word uint32
	}{
		{"str x0,[x1]", "x0", "x1", 0, 0xf9000020},
		{"str w0,[x29,#0xc]", "w0", "x29", 0xc, 0xb9000fa0},
	}
	for _, c := range cases {
		in, err := New().Str(reg(t, c.rt), reg(t, c.rn), c.off)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Str(reg(t, first.rt), reg(t, first.rn), first.off)
	require.NoError(t, err)
	_, ok := in.(Str)
	require.True(t, ok, "type = %T, want Str", in)
}
