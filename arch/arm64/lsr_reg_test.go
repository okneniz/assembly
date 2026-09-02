package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLsrRegCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"lsr x1,x2,x3",
			ctorLsrReg(t, xreg(t, 1), xreg(t, 2), xreg(t, 3)),
			0x9a032441,
		},
		{
			"lsr w1,w2,w3",
			ctorLsrReg(t, wreg(t, 1), wreg(t, 2), wreg(t, 3)),
			0x1a032441,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
	}

	in := ctorLsrReg(t, xreg(t, 1), xreg(t, 2), xreg(t, 3))
	_, ok := in.(LsrReg)
	require.True(t, ok, "type = %T, want LsrReg", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"lsr x+w",
			func() error {
				_, err := New().LsrReg(xreg(t, 0), wreg(t, 1), xreg(t, 2))
				return err
			},
		},
		{
			"lsr sp",
			func() error {
				_, err := New().LsrReg(SP, xreg(t, 1), xreg(t, 2))
				return err
			},
		},
		{
			"lsr rn sp",
			func() error {
				_, err := New().LsrReg(xreg(t, 0), SP, xreg(t, 2))
				return err
			},
		},
		{
			"lsr rm sp",
			func() error {
				_, err := New().LsrReg(xreg(t, 0), xreg(t, 1), SP)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorLsrReg — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorLsrReg(t *testing.T, rd, rn, rm Reg) Instr {
	t.Helper()
	in, err := New().LsrReg(rd, rn, rm)
	require.NoError(t, err)
	return in
}
