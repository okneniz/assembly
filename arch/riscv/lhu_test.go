package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLhuCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"lhu a0, 0(a1)", New().Lhu(xreg(t, 10), xreg(t, 11), off(t, 0)), 0x0005d503},
		{"lhu zero, -2048(ra)", New().Lhu(xreg(t, 0), xreg(t, 1), off(t, -2048)), 0x8000d003},
		{"lhu t6, 2047(sp)", New().Lhu(xreg(t, 31), xreg(t, 2), off(t, 2047)), 0x7ff15f83},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
