package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSraiCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"srai a0, a1, 0", New().Srai(xreg(t, 10), xreg(t, 11), imm12(t, 0)), 0x4005d513},
		{"srai a0, a1, 32", New().Srai(xreg(t, 10), xreg(t, 11), imm12(t, 32)), 0x4205d513},
		{"srai a0, a1, 63", New().Srai(xreg(t, 10), xreg(t, 11), imm12(t, 63)), 0x43f5d513},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}

	// rd == rs1 in x8-x15 with shamt >= 32 compresses to c.srai (2 bytes).
	b := ctorBytes(t, New().Srai(xreg(t, 10), xreg(t, 10), imm12(t, 40)))
	require.Len(t, b, 2, "srai a0,a0,40 (c.srai)")
}
