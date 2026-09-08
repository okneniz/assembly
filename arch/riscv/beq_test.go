package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBeqCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"beq a0, a1, off 0x10", New().Beq(xreg(t, 10), xreg(t, 11), 0x10), 0x00b50863},
		{"beq t0, t1, off -0x8", New().Beq(xreg(t, 5), xreg(t, 6), -0x8), 0xfe628ce3},
		{"beq a0, a1, off 0xffe", New().Beq(xreg(t, 10), xreg(t, 11), 0xffe), 0x7eb50fe3},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}

	// beqz with an x8-x15 register and a small offset compresses to c.beqz (2 bytes).
	b := ctorBytes(t, New().Beq(xreg(t, 10), xreg(t, 0), 0x8))
	require.Len(t, b, 2, "beq a0,zero,0x1008 (c.beqz)")
}
