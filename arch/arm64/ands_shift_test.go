package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAndsShiftBuild(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"tst w1,w2",
			buildAndsShift(t, WZR, wreg(t, 1), wreg(t, 2), imm6(t, 0), LSL),
			0x6a02003f,
		},
		{
			"ands x1,x2,x3,lsl#4",
			buildAndsShift(t, xreg(t, 1), xreg(t, 2), xreg(t, 3), imm6(t, 4), LSL),
			0xea031041,
		},
		{
			"ands w1,w2,w3,asr#5",
			buildAndsShift(t, wreg(t, 1), wreg(t, 2), wreg(t, 3), imm6(t, 5), ASR),
			0x6a831441,
		},
	}
	for _, c := range cases {
		got := buildWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
	}

	in := buildAndsShift(t, xreg(t, 1), xreg(t, 2), xreg(t, 3), imm6(t, 1), LSL)
	_, ok := in.(AndsShift)
	require.True(t, ok, "type = %T, want AndsShift", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"ands x+w",
			func() error {
				_, err := New().AndsShift(xreg(t, 0), wreg(t, 1), xreg(t, 2), imm6(t, 1), LSL)
				return err
			},
		},
		{
			"andsshift sp",
			func() error {
				_, err := New().AndsShift(SP, xreg(t, 1), xreg(t, 2), imm6(t, 1), LSL)
				return err
			},
		},
		{
			"andsshift rn sp",
			func() error {
				_, err := New().AndsShift(xreg(t, 0), SP, xreg(t, 2), imm6(t, 1), LSL)
				return err
			},
		},
		{
			"andsshift rm sp",
			func() error {
				_, err := New().AndsShift(xreg(t, 0), xreg(t, 1), SP, imm6(t, 1), LSL)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// buildAndsShift — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func buildAndsShift(t *testing.T, rd, rn, rm Reg, imm Imm6, sh Shift) Instr {
	t.Helper()
	in, err := New().AndsShift(rd, rn, rm, imm, sh)
	require.NoError(t, err)
	return in
}
