package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddImmCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"add x0,x1,#0x42",
			ctorAddImm(t, xreg(t, 0), xreg(t, 1), imm12(t, 0x42), NoSh12),
			0x91010820,
		},
		{
			"add x0,x1,#1,lsl#12",
			ctorAddImm(t, xreg(t, 0), xreg(t, 1), imm12(t, 1), LSL12),
			0x91400420,
		},
		{
			"add w2,w3,#7",
			ctorAddImm(t, wreg(t, 2), wreg(t, 3), imm12(t, 7), NoSh12),
			0x11001c62,
		},
		{
			"add sp,sp,#0x10",
			ctorAddImm(t, SP, SP, imm12(t, 0x10), NoSh12),
			0x910043ff,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
	}

	in := ctorAddImm(t, xreg(t, 0), xreg(t, 1), imm12(t, 1), NoSh12)
	_, ok := in.(AddImm)
	require.True(t, ok, "type = %T, want AddImm", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"add xzr",
			func() error {
				_, err := New().AddImm(XZR, xreg(t, 1), imm12(t, 1), NoSh12)
				return err
			},
		},
		{
			"add rn xzr",
			func() error {
				_, err := New().AddImm(xreg(t, 0), XZR, imm12(t, 1), NoSh12)
				return err
			},
		},
		{
			"add sp+w",
			func() error {
				_, err := New().AddImm(SP, wreg(t, 1), imm12(t, 1), NoSh12)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}
