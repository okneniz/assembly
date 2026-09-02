package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEonShiftCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"eon x1,x2,x3,ror#5",
			ctorEonShift(t, xreg(t, 1), xreg(t, 2), xreg(t, 3), imm6(t, 5), ROR),
			0xcae31441,
		},
		{
			"eon w1,w2,w3",
			ctorEonShift(t, wreg(t, 1), wreg(t, 2), wreg(t, 3), imm6(t, 0), LSL),
			0x4a230041,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
	}

	in := ctorEonShift(t, xreg(t, 1), xreg(t, 2), xreg(t, 3), imm6(t, 1), LSL)
	_, ok := in.(EonShift)
	require.True(t, ok, "type = %T, want EonShift", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"eon x+w",
			func() error {
				_, err := New().EonShift(xreg(t, 0), wreg(t, 1), xreg(t, 2), imm6(t, 1), LSL)
				return err
			},
		},
		{
			"eonshift sp",
			func() error {
				_, err := New().EonShift(SP, xreg(t, 1), xreg(t, 2), imm6(t, 1), LSL)
				return err
			},
		},
		{
			"eonshift rn sp",
			func() error {
				_, err := New().EonShift(xreg(t, 0), SP, xreg(t, 2), imm6(t, 1), LSL)
				return err
			},
		},
		{
			"eonshift rm sp",
			func() error {
				_, err := New().EonShift(xreg(t, 0), xreg(t, 1), SP, imm6(t, 1), LSL)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorEonShift — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorEonShift(t *testing.T, rd, rn, rm Reg, imm Imm6, sh Shift) Instr {
	t.Helper()
	in, err := New().EonShift(rd, rn, rm, imm, sh)
	require.NoError(t, err)
	return in
}
