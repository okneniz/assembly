package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddiCtor(t *testing.T) {
	require.Equal(
		t,
		uint32(0xfff58513),
		ctorWord(t, NewAddi(xreg(t, 10), xreg(t, 11), imm12(t, -1))),
	)
	// rs1 = zero is printed as li and compresses to c.li (2 bytes).
	b := ctorBytes(t, NewAddi(xreg(t, 10), Zero, imm12(t, 5)))
	require.Len(t, b, 2, "li a0,5 (c.li)")
	in := NewAddi(xreg(t, 1), xreg(t, 2), imm12(t, 0))
	_, ok := in.(Addi)
	require.True(t, ok, "type = %T, want Addi", in)
	// Imm outside -2048..2047 is rejected by the Imm12 constructor.
	_, err := NewImm12(2048)
	require.Error(t, err)
	_, err = NewImm12(-2049)
	require.Error(t, err)
}
