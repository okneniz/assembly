package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMovzCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"movz x0,#0x1234",
			ctorMovz(t, xreg(t, 0), imm16(t, 0x1234), Hw0),
			0xd2824680,
		},
		{
			"movz x0,#1",
			ctorMovz(t, xreg(t, 0), imm16(t, 1), Hw0),
			0xd2800020,
		},
		{
			"movz w3,#0x42,hw1",
			ctorMovz(t, wreg(t, 3), imm16(t, 0x42), Hw1),
			0x52a00843,
		},
		{
			"movz xzr,#0,hw3",
			ctorMovz(t, XZR, imm16(t, 0), Hw3),
			0xd2e0001f,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
	}

	in := ctorMovz(t, xreg(t, 0), imm16(t, 1), Hw0)
	_, ok := in.(Movz)
	require.True(t, ok, "type = %T, want Movz", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"movz sp",
			func() error {
				_, err := New().Movz(SP, imm16(t, 1), Hw0)
				return err
			},
		},
		{
			"movz w0,hw2",
			func() error {
				_, err := New().Movz(wreg(t, 0), imm16(t, 1), Hw2)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}
