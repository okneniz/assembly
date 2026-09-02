package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSubCtor(t *testing.T) {
	require.Equal(
		t,
		uint32(0x407302b3),
		ctorWord(t, New().Sub(xreg(t, 5), xreg(t, 6), xreg(t, 7))),
	)
	in := New().Sub(xreg(t, 1), xreg(t, 2), xreg(t, 3))
	_, ok := in.(Sub)
	require.True(t, ok, "type = %T, want Sub", in)
}
