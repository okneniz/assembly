package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLbuCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"lbu a0, 0(a1)", New().Lbu(xreg(t, 10), xreg(t, 11), off(t, 0)), 0x0005c503},
		{"lbu zero, -2048(ra)", New().Lbu(xreg(t, 0), xreg(t, 1), off(t, -2048)), 0x8000c003},
		{"lbu t6, 2047(sp)", New().Lbu(xreg(t, 31), xreg(t, 2), off(t, 2047)), 0x7ff14f83},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
