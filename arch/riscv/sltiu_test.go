package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSltiuCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"sltiu a0, a1, 2047", New().Sltiu(xreg(t, 10), xreg(t, 11), imm12(t, 2047)), 0x7ff5b513},
		{"sltiu a0, a1, -2048", New().Sltiu(xreg(t, 10), xreg(t, 11), imm12(t, -2048)), 0x8005b513},
		{"sltiu a0, a1, 1", New().Sltiu(xreg(t, 10), xreg(t, 11), imm12(t, 1)), 0x0015b513},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
