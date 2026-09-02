package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSltiCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"slti a0, a1, 2047", New().Slti(xreg(t, 10), xreg(t, 11), imm12(t, 2047)), 0x7ff5a513},
		{"slti a0, a1, -2048", New().Slti(xreg(t, 10), xreg(t, 11), imm12(t, -2048)), 0x8005a513},
		{"slti a0, a1, 0", New().Slti(xreg(t, 10), xreg(t, 11), imm12(t, 0)), 0x0005a513},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
