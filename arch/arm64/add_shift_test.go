package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddShiftCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"add x1,x2,x3",
			ctorAddShift(t, xreg(t, 1), xreg(t, 2), xreg(t, 3), imm6(t, 0), LSL),
			0x8b030041,
		},
		{
			"add x1,x2,x3,lsl#4",
			ctorAddShift(t, xreg(t, 1), xreg(t, 2), xreg(t, 3), imm6(t, 4), LSL),
			0x8b031041,
		},
		{
			"add w1,w2,w3,asr#5",
			ctorAddShift(t, wreg(t, 1), wreg(t, 2), wreg(t, 3), imm6(t, 5), ASR),
			0x0b831441,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
	}

	in := ctorAddShift(t, xreg(t, 1), xreg(t, 2), xreg(t, 3), imm6(t, 1), LSL)
	_, ok := in.(AddShift)
	require.True(t, ok, "type = %T, want AddShift", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"add x+w",
			func() error {
				_, err := New().AddShift(xreg(t, 0), wreg(t, 1), xreg(t, 2), imm6(t, 1), LSL)
				return err
			},
		},
		{
			"addshift sp",
			func() error {
				_, err := New().AddShift(SP, xreg(t, 1), xreg(t, 2), imm6(t, 1), LSL)
				return err
			},
		},
		{
			"add ror",
			func() error {
				_, err := New().AddShift(xreg(t, 0), xreg(t, 1), xreg(t, 2), imm6(t, 1), ROR)
				return err
			},
		},
		{
			"add w + imm6=32",
			func() error {
				_, err := New().AddShift(wreg(t, 0), wreg(t, 1), wreg(t, 2), imm6(t, 32), LSL)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}
