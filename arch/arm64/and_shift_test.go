package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAndShiftCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"and x1,x2,x3",
			ctorAndShift(t, xreg(t, 1), xreg(t, 2), xreg(t, 3), imm6(t, 0), LSL),
			0x8a030041,
		},
		{
			"and x1,x2,x3,lsl#63",
			ctorAndShift(t, xreg(t, 1), xreg(t, 2), xreg(t, 3), imm6(t, 63), LSL),
			0x8a03fc41,
		},
		{
			"and w1,w2,w3,ror#5",
			ctorAndShift(t, wreg(t, 1), wreg(t, 2), wreg(t, 3), imm6(t, 5), ROR),
			0x0ac31441,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
	}

	in := ctorAndShift(t, xreg(t, 1), xreg(t, 2), xreg(t, 3), imm6(t, 1), LSL)
	_, ok := in.(AndShift)
	require.True(t, ok, "type = %T, want AndShift", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"and x+w",
			func() error {
				_, err := New().AndShift(xreg(t, 0), wreg(t, 1), xreg(t, 2), imm6(t, 1), LSL)
				return err
			},
		},
		{
			"andshift sp",
			func() error {
				_, err := New().AndShift(SP, xreg(t, 1), xreg(t, 2), imm6(t, 1), LSL)
				return err
			},
		},
		{
			"andshift rn sp",
			func() error {
				_, err := New().AndShift(xreg(t, 0), SP, xreg(t, 2), imm6(t, 1), LSL)
				return err
			},
		},
		{
			"andshift rm sp",
			func() error {
				_, err := New().AndShift(xreg(t, 0), xreg(t, 1), SP, imm6(t, 1), LSL)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorAndShift — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorAndShift(t *testing.T, rd, rn, rm Reg, imm Imm6, sh Shift) Instr {
	t.Helper()
	in, err := New().AndShift(rd, rn, rm, imm, sh)
	require.NoError(t, err)
	return in
}
