package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRetCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"ret",
			ctorRet(t, xreg(t, 30)),
			0xd65f03c0,
		},
		{
			"ret x8",
			ctorRet(t, xreg(t, 8)),
			0xd65f0100,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
	}

	in := ctorRet(t, xreg(t, 30))
	_, ok := in.(Ret)
	require.True(t, ok, "type = %T, want Ret", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"ret sp",
			func() error {
				_, err := New().Ret(SP)
				return err
			},
		},
		{
			"ret w0",
			func() error {
				_, err := New().Ret(wreg(t, 0))
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}
