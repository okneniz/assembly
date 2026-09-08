package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBneCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"bne a0, a1, off 0x10", New().Bne(xreg(t, 10), xreg(t, 11), 0x10), 0x00b51863},
		{"bne t0, t1, off -0x8", New().Bne(xreg(t, 5), xreg(t, 6), -0x8), 0xfe629ce3},
		{"bne a0, a1, off 0xffe", New().Bne(xreg(t, 10), xreg(t, 11), 0xffe), 0x7eb51fe3},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}

	// bnez with an x8-x15 register and a small offset compresses to c.bnez (2 bytes).
	b := ctorBytes(t, New().Bne(xreg(t, 10), xreg(t, 0), 0x8))
	require.Len(t, b, 2, "bne a0,zero,0x1008 (c.bnez)")
}
