package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOrnShiftCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"orn x1,x2,x3,lsr#5",
			ctorOrnShift(t, xreg(t, 1), xreg(t, 2), xreg(t, 3), imm6(t, 5), LSR),
			0xaa631441,
		},
		{
			"mvn x1,x2",
			ctorOrnShift(t, xreg(t, 1), XZR, xreg(t, 2), imm6(t, 0), LSL),
			0xaa2203e1,
		},
		{
			"orn w1,w2,w3",
			ctorOrnShift(t, wreg(t, 1), wreg(t, 2), wreg(t, 3), imm6(t, 0), LSL),
			0x2a230041,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
	}

	in := ctorOrnShift(t, xreg(t, 1), xreg(t, 2), xreg(t, 3), imm6(t, 1), LSL)
	_, ok := in.(OrnShift)
	require.True(t, ok, "type = %T, want OrnShift", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"orn x+w",
			func() error {
				_, err := New().OrnShift(xreg(t, 0), wreg(t, 1), xreg(t, 2), imm6(t, 1), LSL)
				return err
			},
		},
		{
			"ornshift sp",
			func() error {
				_, err := New().OrnShift(xreg(t, 0), SP, xreg(t, 2), imm6(t, 1), LSL)
				return err
			},
		},
		{
			"ornshift rd sp",
			func() error {
				_, err := New().OrnShift(SP, xreg(t, 1), xreg(t, 2), imm6(t, 1), LSL)
				return err
			},
		},
		{
			"ornshift rm sp",
			func() error {
				_, err := New().OrnShift(xreg(t, 0), xreg(t, 1), SP, imm6(t, 1), LSL)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorOrnShift — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorOrnShift(t *testing.T, rd, rn, rm Reg, imm Imm6, sh Shift) Instr {
	t.Helper()
	in, err := New().OrnShift(rd, rn, rm, imm, sh)
	require.NoError(t, err)
	return in
}
