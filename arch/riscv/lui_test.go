package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLuiCtor(t *testing.T) {
	require.Equal(
		t,
		uint32(0x123452b7),
		ctorWord(t, New().Lui(xreg(t, 5), imm20(t, 0x12345))),
	)
	in := New().Lui(xreg(t, 1), imm20(t, 1))
	_, ok := in.(Lui)
	require.True(t, ok, "type = %T, want Lui", in)
	// Imm20 outside 0..0xfffff is rejected by the Imm20 constructor.
	_, err := New().Imm20(-1)
	require.Error(t, err)
	_, err = New().Imm20(0x100000)
	require.Error(t, err)
}
