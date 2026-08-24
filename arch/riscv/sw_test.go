package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSwCtor(t *testing.T) {
	b := ctorBytes(t, NewSw(xreg(t, 10), Sp, off(t, 12)))
	require.Len(t, b, 2, "sw a0,12(sp) (c.swsp)")
	in := NewSw(xreg(t, 1), xreg(t, 2), off(t, 0))
	_, ok := in.(Sw)
	require.True(t, ok, "type = %T, want Sw", in)
}
