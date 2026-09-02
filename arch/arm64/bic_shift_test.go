package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBicShiftCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"bic x1,x2,x3",
			ctorBicShift(t, xreg(t, 1), xreg(t, 2), xreg(t, 3), imm6(t, 0), LSL),
			0x8a230041,
		},
		{
			"bic w1,w2,w3,asr#5",
			ctorBicShift(t, wreg(t, 1), wreg(t, 2), wreg(t, 3), imm6(t, 5), ASR),
			0x0aa31441,
		},
		{
			"bic x1,x2,x3,lsr#63",
			ctorBicShift(t, xreg(t, 1), xreg(t, 2), xreg(t, 3), imm6(t, 63), LSR),
			0x8a63fc41,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
	}

	in := ctorBicShift(t, xreg(t, 1), xreg(t, 2), xreg(t, 3), imm6(t, 1), LSL)
	_, ok := in.(BicShift)
	require.True(t, ok, "type = %T, want BicShift", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"bic x+w",
			func() error {
				_, err := New().BicShift(xreg(t, 0), wreg(t, 1), xreg(t, 2), imm6(t, 1), LSL)
				return err
			},
		},
		{
			"bicshift sp",
			func() error {
				_, err := New().BicShift(SP, xreg(t, 1), xreg(t, 2), imm6(t, 1), LSL)
				return err
			},
		},
		{
			"bicshift rn sp",
			func() error {
				_, err := New().BicShift(xreg(t, 0), SP, xreg(t, 2), imm6(t, 1), LSL)
				return err
			},
		},
		{
			"bicshift rm sp",
			func() error {
				_, err := New().BicShift(xreg(t, 0), xreg(t, 1), SP, imm6(t, 1), LSL)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorBicShift — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorBicShift(t *testing.T, rd, rn, rm Reg, imm Imm6, sh Shift) Instr {
	t.Helper()
	in, err := New().BicShift(rd, rn, rm, imm, sh)
	require.NoError(t, err)
	return in
}
