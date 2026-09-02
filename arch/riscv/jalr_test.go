package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJalrCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"jalr t0, 8(a0)", New().Jalr(xreg(t, 5), xreg(t, 10), off(t, 8)), 0x008502e7},
		{"jalr ra, 2047(a0)", New().Jalr(xreg(t, 1), xreg(t, 10), off(t, 2047)), 0x7ff500e7},
		{"jalr t0, -2048(a1)", New().Jalr(xreg(t, 5), xreg(t, 11), off(t, -2048)), 0x800582e7},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}

	// off = 0 with rd = zero, rs1 = ra compresses to c.jr (2 bytes).
	b := ctorBytes(t, New().Jalr(xreg(t, 0), xreg(t, 1), off(t, 0)))
	require.Len(t, b, 2, "jalr zero,0(ra) (c.jr)")
}
