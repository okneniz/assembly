package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSlliwCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"slliw a0, a1, 0", New().Slliw(xreg(t, 10), xreg(t, 11), imm12(t, 0)), 0x0005951b},
		{"slliw a0, a1, 1", New().Slliw(xreg(t, 10), xreg(t, 11), imm12(t, 1)), 0x0015951b},
		{"slliw a0, a1, 31", New().Slliw(xreg(t, 10), xreg(t, 11), imm12(t, 31)), 0x01f5951b},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
