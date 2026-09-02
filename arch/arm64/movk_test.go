package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMovkCtor(t *testing.T) {
	got := ctorWord(t, ctorMovk(t, xreg(t, 0), imm16(t, 0x1234), Hw1))
	require.Equal(t, uint32(0xf2a24680), got, "movk")
	got = ctorWord(t, ctorMovk(t, wreg(t, 4), imm16(t, 0x1234), Hw0))
	require.Equal(t, uint32(0x72824684), got, "movk w4")
	in := ctorMovk(t, xreg(t, 0), imm16(t, 1), Hw0)
	_, ok := in.(Movk)
	require.True(t, ok, "type = %T, want Movk", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"movk w0,hw3",
			func() error {
				_, err := New().Movk(wreg(t, 0), imm16(t, 1), Hw3)
				return err
			},
		},
		{
			"movk wsp",
			func() error {
				_, err := New().Movk(WSP, imm16(t, 1), Hw0)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}
