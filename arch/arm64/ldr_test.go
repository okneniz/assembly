package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLdrCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"ldr x0,[x1]",
			ctorLdr(t, xreg(t, 0), xreg(t, 1), 0),
			0xf9400020,
		},
		{
			"ldr x0,[x1,#8]",
			ctorLdr(t, xreg(t, 0), xreg(t, 1), 8),
			0xf9400420,
		},
		{
			"ldr w2,[sp,#0x10]",
			ctorLdr(t, wreg(t, 2), SP, 0x10),
			0xb94013e2,
		},
		{
			"ldr xzr,[x1,#0x7ff8]",
			ctorLdr(t, XZR, xreg(t, 1), 0x7ff8),
			0xf97ffc3f,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
	}

	in := ctorLdr(t, xreg(t, 0), xreg(t, 1), 0)
	_, ok := in.(Ldr)
	require.True(t, ok, "type = %T, want Ldr", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"ldr sp,rt",
			func() error {
				_, err := NewLdr(SP, xreg(t, 1), 0)
				return err
			},
		},
		{
			"ldr w,[w1]",
			func() error {
				_, err := NewLdr(wreg(t, 0), wreg(t, 1), 0)
				return err
			},
		},
		{
			"ldr xzr base",
			func() error {
				_, err := NewLdr(xreg(t, 0), XZR, 0)
				return err
			},
		},
		{
			"ldr off 4 in x form",
			func() error {
				_, err := NewLdr(xreg(t, 0), xreg(t, 1), 4)
				return err
			},
		},
		{
			"ldr off 0x8000 in w form",
			func() error {
				_, err := NewLdr(wreg(t, 0), xreg(t, 1), 0x8000)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}
