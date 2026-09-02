package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddiwCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"addiw a0, a1, 2047", New().Addiw(xreg(t, 10), xreg(t, 11), imm12(t, 2047)), 0x7ff5851b},
		{"addiw a0, a1, -2048", New().Addiw(xreg(t, 10), xreg(t, 11), imm12(t, -2048)), 0x8005851b},
		{"addiw a0, a1, 16", New().Addiw(xreg(t, 10), xreg(t, 11), imm12(t, 16)), 0x0105851b},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}

	// rd == rs1 with a 6-bit imm compresses to c.addiw (2 bytes).
	b := ctorBytes(t, New().Addiw(xreg(t, 10), xreg(t, 10), imm12(t, -1)))
	require.Len(t, b, 2, "addiw a0,a0,-1 (c.addiw)")
}
