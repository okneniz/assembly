package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddCtor(t *testing.T) {
	require.Equal(
		t,
		uint32(0x00c58533),
		ctorWord(t, NewAdd(xreg(t, 10), xreg(t, 11), xreg(t, 12))),
	)
	// rd == rs1 compresses to c.add (2 bytes), decode agrees.
	b := ctorBytes(t, NewAdd(xreg(t, 10), xreg(t, 10), xreg(t, 12)))
	require.Len(t, b, 2, "add a0,a0,a2 (c.add)")
	in := NewAdd(xreg(t, 1), xreg(t, 2), xreg(t, 3))
	_, ok := in.(Add)
	require.True(t, ok, "type = %T, want Add", in)
	// Register number outside 0..31 is rejected by the Reg constructor.
	_, err := X(32)
	require.Error(t, err)
}
