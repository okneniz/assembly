package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtrCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"extr x1,x2,x3,#5",
			ctorExtr(t, xreg(t, 1), xreg(t, 2), xreg(t, 3), imm6(t, 5)),
			0x93031441,
		},
		{
			"ror x1,x2,#4",
			ctorExtr(t, xreg(t, 1), xreg(t, 2), xreg(t, 2), imm6(t, 4)),
			0x93021041,
		},
		{
			"extr x1,x2,x3,#63",
			ctorExtr(t, xreg(t, 1), xreg(t, 2), xreg(t, 3), imm6(t, 63)),
			0x9303fc41,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
	}

	in := ctorExtr(t, xreg(t, 1), xreg(t, 2), xreg(t, 3), imm6(t, 5))
	_, ok := in.(Extr)
	require.True(t, ok, "type = %T, want Extr", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"extr w form",
			func() error {
				_, err := New().Extr(wreg(t, 1), xreg(t, 2), xreg(t, 3), imm6(t, 5))
				return err
			},
		},
		{
			"extr sp",
			func() error {
				_, err := New().Extr(xreg(t, 1), SP, xreg(t, 3), imm6(t, 5))
				return err
			},
		},
		{
			"extr rm sp",
			func() error {
				_, err := New().Extr(xreg(t, 1), xreg(t, 2), SP, imm6(t, 5))
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorExtr — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorExtr(t *testing.T, rd, rn, rm Reg, lsb Imm6) Instr {
	t.Helper()
	in, err := New().Extr(rd, rn, rm, lsb)
	require.NoError(t, err)
	return in
}
