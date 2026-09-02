package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOriCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"ori a0, a1, 2047", New().Ori(xreg(t, 10), xreg(t, 11), imm12(t, 2047)), 0x7ff5e513},
		{"ori a0, a1, -2048", New().Ori(xreg(t, 10), xreg(t, 11), imm12(t, -2048)), 0x8005e513},
		{"ori a0, a1, 0x7f0", New().Ori(xreg(t, 10), xreg(t, 11), imm12(t, 0x7f0)), 0x7f05e513},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
