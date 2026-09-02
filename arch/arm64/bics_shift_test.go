package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBicsShiftCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"bics x1,x2,x3",
			ctorBicsShift(t, xreg(t, 1), xreg(t, 2), xreg(t, 3), imm6(t, 0), LSL),
			0xea230041,
		},
		{
			"bics w1,w2,w3,lsl#4",
			ctorBicsShift(t, wreg(t, 1), wreg(t, 2), wreg(t, 3), imm6(t, 4), LSL),
			0x6a231041,
		},
		{
			"bics x1,xzr,x3,ror#2",
			ctorBicsShift(t, xreg(t, 1), XZR, xreg(t, 3), imm6(t, 2), ROR),
			0xeae30be1,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
	}

	in := ctorBicsShift(t, xreg(t, 1), xreg(t, 2), xreg(t, 3), imm6(t, 1), LSL)
	_, ok := in.(BicsShift)
	require.True(t, ok, "type = %T, want BicsShift", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"bics x+w",
			func() error {
				_, err := New().BicsShift(xreg(t, 0), wreg(t, 1), xreg(t, 2), imm6(t, 1), LSL)
				return err
			},
		},
		{
			"bicsshift sp",
			func() error {
				_, err := New().BicsShift(SP, xreg(t, 1), xreg(t, 2), imm6(t, 1), LSL)
				return err
			},
		},
		{
			"bicsshift rn sp",
			func() error {
				_, err := New().BicsShift(xreg(t, 0), SP, xreg(t, 2), imm6(t, 1), LSL)
				return err
			},
		},
		{
			"bicsshift rm sp",
			func() error {
				_, err := New().BicsShift(xreg(t, 0), xreg(t, 1), SP, imm6(t, 1), LSL)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorBicsShift — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorBicsShift(t *testing.T, rd, rn, rm Reg, imm Imm6, sh Shift) Instr {
	t.Helper()
	in, err := New().BicsShift(rd, rn, rm, imm, sh)
	require.NoError(t, err)
	return in
}
