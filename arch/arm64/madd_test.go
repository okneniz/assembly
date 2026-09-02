package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMaddCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"madd x1,x2,x3,x4",
			ctorMadd(t, xreg(t, 1), xreg(t, 2), xreg(t, 3), xreg(t, 4)),
			0x9b031041,
		},
		{
			"mul x1,x2,x3",
			ctorMadd(t, xreg(t, 1), xreg(t, 2), xreg(t, 3), XZR),
			0x9b037c41,
		},
		{
			"madd w1,w2,w3,w4",
			ctorMadd(t, wreg(t, 1), wreg(t, 2), wreg(t, 3), wreg(t, 4)),
			0x1b031041,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
	}

	in := ctorMadd(t, xreg(t, 1), xreg(t, 2), xreg(t, 3), xreg(t, 4))
	_, ok := in.(Madd)
	require.True(t, ok, "type = %T, want Madd", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"madd x+w",
			func() error {
				_, err := New().Madd(xreg(t, 1), wreg(t, 2), xreg(t, 3), xreg(t, 4))
				return err
			},
		},
		{
			"madd sp",
			func() error {
				_, err := New().Madd(xreg(t, 1), xreg(t, 2), xreg(t, 3), SP)
				return err
			},
		},
		{
			"madd rd sp",
			func() error {
				_, err := New().Madd(SP, xreg(t, 2), xreg(t, 3), xreg(t, 4))
				return err
			},
		},
		{
			"madd rn sp",
			func() error {
				_, err := New().Madd(xreg(t, 1), SP, xreg(t, 3), xreg(t, 4))
				return err
			},
		},
		{
			"madd rm sp",
			func() error {
				_, err := New().Madd(xreg(t, 1), xreg(t, 2), SP, xreg(t, 4))
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorMadd — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorMadd(t *testing.T, rd, rn, rm, ra Reg) Instr {
	t.Helper()
	in, err := New().Madd(rd, rn, rm, ra)
	require.NoError(t, err)
	return in
}
