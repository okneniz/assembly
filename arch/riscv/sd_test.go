package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSdCtor(t *testing.T) {
	b := ctorBytes(t, New().Sd(xreg(t, 10), Sp, off(t, 0)))
	require.Len(t, b, 2, "sd a0,0(sp) (c.sdsp)")
	in := New().Sd(xreg(t, 1), xreg(t, 2), off(t, 0))
	_, ok := in.(Sd)
	require.True(t, ok, "type = %T, want Sd", in)
}
