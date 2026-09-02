package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSlliCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"slli a0, a1, 0", New().Slli(xreg(t, 10), xreg(t, 11), imm12(t, 0)), 0x00059513},
		{"slli a0, a1, 1", New().Slli(xreg(t, 10), xreg(t, 11), imm12(t, 1)), 0x00159513},
		{"slli a0, a1, 63", New().Slli(xreg(t, 10), xreg(t, 11), imm12(t, 63)), 0x03f59513},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}

	// rd == rs1 with 0 < shamt < 64 compresses to c.slli (2 bytes).
	b := ctorBytes(t, New().Slli(xreg(t, 10), xreg(t, 10), imm12(t, 4)))
	require.Len(t, b, 2, "slli a0,a0,4 (c.slli)")
}
