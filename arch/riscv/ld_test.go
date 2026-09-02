package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLdCtor(t *testing.T) {
	require.Equal(
		t,
		uint32(0xff813503),
		ctorWord(t, New().Ld(xreg(t, 10), Sp, off(t, -8))),
	)
	in := New().Ld(xreg(t, 1), xreg(t, 2), off(t, 0))
	_, ok := in.(Ld)
	require.True(t, ok, "type = %T, want Ld", in)
	// Off outside -2048..2047 is rejected by the Off constructor.
	_, err := New().Off(-2049)
	require.Error(t, err)
}
