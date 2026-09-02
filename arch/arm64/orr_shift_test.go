package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOrrShiftCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"orr x1,x2,x3,lsl#4",
			ctorOrrShift(t, xreg(t, 1), xreg(t, 2), xreg(t, 3), imm6(t, 4), LSL),
			0xaa031041,
		},
		{
			"orr w1,w2,w3",
			ctorOrrShift(t, wreg(t, 1), wreg(t, 2), wreg(t, 3), imm6(t, 0), LSL),
			0x2a030041,
		},
		{
			"mov x1,x3,ror#2",
			ctorOrrShift(t, xreg(t, 1), XZR, xreg(t, 3), imm6(t, 2), ROR),
			0xaac30be1,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
	}

	in := ctorOrrShift(t, xreg(t, 1), xreg(t, 2), xreg(t, 3), imm6(t, 1), LSL)
	_, ok := in.(OrrShift)
	require.True(t, ok, "type = %T, want OrrShift", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"orr x+w",
			func() error {
				_, err := New().OrrShift(xreg(t, 0), wreg(t, 1), xreg(t, 2), imm6(t, 1), LSL)
				return err
			},
		},
		{
			"orrshift sp",
			func() error {
				_, err := New().OrrShift(xreg(t, 0), xreg(t, 1), SP, imm6(t, 1), LSL)
				return err
			},
		},
		{
			"orrshift rd sp",
			func() error {
				_, err := New().OrrShift(SP, xreg(t, 1), xreg(t, 2), imm6(t, 1), LSL)
				return err
			},
		},
		{
			"orrshift rn sp",
			func() error {
				_, err := New().OrrShift(xreg(t, 0), SP, xreg(t, 2), imm6(t, 1), LSL)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorOrrShift — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorOrrShift(t *testing.T, rd, rn, rm Reg, imm Imm6, sh Shift) Instr {
	t.Helper()
	in, err := New().OrrShift(rd, rn, rm, imm, sh)
	require.NoError(t, err)
	return in
}
