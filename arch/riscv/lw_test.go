package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLwCtor(t *testing.T) {
	// small aligned offsets with x-registers compress (c.lw).
	b := ctorBytes(t, NewLw(xreg(t, 10), xreg(t, 11), off(t, 8)))
	require.Len(t, b, 2, "lw a0,8(a1) (c.lw)")
	in := NewLw(xreg(t, 1), xreg(t, 2), off(t, 0))
	_, ok := in.(Lw)
	require.True(t, ok, "type = %T, want Lw", in)
	// Off outside -2048..2047 is rejected by the Off constructor.
	_, err := NewOff(2048)
	require.Error(t, err)
}
