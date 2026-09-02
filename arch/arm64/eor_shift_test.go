package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEorShiftCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"eor x1,x2,x3,lsl#4",
			ctorEorShift(t, xreg(t, 1), xreg(t, 2), xreg(t, 3), imm6(t, 4), LSL),
			0xca031041,
		},
		{
			"eor w1,w2,w3",
			ctorEorShift(t, wreg(t, 1), wreg(t, 2), wreg(t, 3), imm6(t, 0), LSL),
			0x4a030041,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
	}

	in := ctorEorShift(t, xreg(t, 1), xreg(t, 2), xreg(t, 3), imm6(t, 1), LSL)
	_, ok := in.(EorShift)
	require.True(t, ok, "type = %T, want EorShift", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"eor x+w",
			func() error {
				_, err := New().EorShift(xreg(t, 0), wreg(t, 1), xreg(t, 2), imm6(t, 1), LSL)
				return err
			},
		},
		{
			"eorshift sp",
			func() error {
				_, err := New().EorShift(xreg(t, 0), SP, xreg(t, 2), imm6(t, 1), LSL)
				return err
			},
		},
		{
			"eorshift rd sp",
			func() error {
				_, err := New().EorShift(SP, xreg(t, 1), xreg(t, 2), imm6(t, 1), LSL)
				return err
			},
		},
		{
			"eorshift rm sp",
			func() error {
				_, err := New().EorShift(xreg(t, 0), xreg(t, 1), SP, imm6(t, 1), LSL)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorEorShift — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorEorShift(t *testing.T, rd, rn, rm Reg, imm Imm6, sh Shift) Instr {
	t.Helper()
	in, err := New().EorShift(rd, rn, rm, imm, sh)
	require.NoError(t, err)
	return in
}
