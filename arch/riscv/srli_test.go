package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSrliCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"srli a0, a1, 0", New().Srli(xreg(t, 10), xreg(t, 11), imm12(t, 0)), 0x0005d513},
		{"srli a0, a1, 1", New().Srli(xreg(t, 10), xreg(t, 11), imm12(t, 1)), 0x0015d513},
		{"srli a0, a1, 31", New().Srli(xreg(t, 10), xreg(t, 11), imm12(t, 31)), 0x01f5d513},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}

	// rd == rs1 in x8-x15 with 0 < shamt < 32 compresses to c.srli (2 bytes).
	b := ctorBytes(t, New().Srli(xreg(t, 10), xreg(t, 10), imm12(t, 8)))
	require.Len(t, b, 2, "srli a0,a0,8 (c.srli)")
}
