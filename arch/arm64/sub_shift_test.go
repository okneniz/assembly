package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSubShiftCtor(t *testing.T) {
	got := ctorWord(
		t,
		ctorSubShift(t, xreg(t, 1), xreg(t, 2), xreg(t, 3), imm6(t, 2), LSR),
	)
	require.Equal(t, uint32(0xcb430841), got, "sub lsr#2")
	in := ctorSubShift(t, xreg(t, 0), xreg(t, 1), xreg(t, 2), imm6(t, 1), LSL)
	_, ok := in.(SubShift)
	require.True(t, ok, "type = %T, want SubShift", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"subshift sp",
			func() error {
				_, err := NewSubShift(xreg(t, 0), xreg(t, 1), SP, imm6(t, 1), LSL)
				return err
			},
		},
		{
			"subshift w + imm6=32",
			func() error {
				_, err := NewSubShift(wreg(t, 0), wreg(t, 1), wreg(t, 2), imm6(t, 32), LSL)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}
