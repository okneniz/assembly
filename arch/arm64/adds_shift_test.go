package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddsShiftCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"cmn x1,x2",
			ctorAddsShift(t, XZR, xreg(t, 1), xreg(t, 2), imm6(t, 0), LSL),
			0xab02003f,
		},
		{
			"adds x1,x2,x3,lsl#4",
			ctorAddsShift(t, xreg(t, 1), xreg(t, 2), xreg(t, 3), imm6(t, 4), LSL),
			0xab031041,
		},
		{
			"adds w1,w2,w3,asr#5",
			ctorAddsShift(t, wreg(t, 1), wreg(t, 2), wreg(t, 3), imm6(t, 5), ASR),
			0x2b831441,
		},
		{
			"adds x1,x2,x3,lsr#63",
			ctorAddsShift(t, xreg(t, 1), xreg(t, 2), xreg(t, 3), imm6(t, 63), LSR),
			0xab43fc41,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
	}

	in := ctorAddsShift(t, xreg(t, 1), xreg(t, 2), xreg(t, 3), imm6(t, 1), LSL)
	_, ok := in.(AddsShift)
	require.True(t, ok, "type = %T, want AddsShift", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"adds x+w",
			func() error {
				_, err := New().AddsShift(xreg(t, 0), wreg(t, 1), xreg(t, 2), imm6(t, 1), LSL)
				return err
			},
		},
		{
			"addsshift sp",
			func() error {
				_, err := New().AddsShift(SP, xreg(t, 1), xreg(t, 2), imm6(t, 1), LSL)
				return err
			},
		},
		{
			"addsshift rn sp",
			func() error {
				_, err := New().AddsShift(xreg(t, 0), SP, xreg(t, 2), imm6(t, 1), LSL)
				return err
			},
		},
		{
			"addsshift rm sp",
			func() error {
				_, err := New().AddsShift(xreg(t, 0), xreg(t, 1), SP, imm6(t, 1), LSL)
				return err
			},
		},
		{
			"adds ror",
			func() error {
				_, err := New().AddsShift(xreg(t, 0), xreg(t, 1), xreg(t, 2), imm6(t, 1), ROR)
				return err
			},
		},
		{
			"adds w + imm6=32",
			func() error {
				_, err := New().AddsShift(wreg(t, 0), wreg(t, 1), wreg(t, 2), imm6(t, 32), LSL)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorAddsShift — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorAddsShift(t *testing.T, rd, rn, rm Reg, imm Imm6, sh Shift) Instr {
	t.Helper()
	in, err := New().AddsShift(rd, rn, rm, imm, sh)
	require.NoError(t, err)
	return in
}
