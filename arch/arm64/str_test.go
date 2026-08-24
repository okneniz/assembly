package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStrCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"str x0,[x1]",
			ctorStr(t, xreg(t, 0), xreg(t, 1), 0),
			0xf9000020,
		},
		{
			"str w0,[x29,#0xc]",
			ctorStr(t, wreg(t, 0), xreg(t, 29), 0xc),
			0xb9000fa0,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
	}

	in := ctorStr(t, xreg(t, 0), xreg(t, 1), 0)
	_, ok := in.(Str)
	require.True(t, ok, "type = %T, want Str", in)
	_, err := NewStr(xreg(t, 0), xreg(t, 1), -8)
	assertErr(t, "str negative offset", err)
}
