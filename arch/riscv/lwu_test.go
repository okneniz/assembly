package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLwuCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"lwu a0, 0(a1)", New().Lwu(xreg(t, 10), xreg(t, 11), off(t, 0)), 0x0005e503},
		{"lwu zero, -2048(ra)", New().Lwu(xreg(t, 0), xreg(t, 1), off(t, -2048)), 0x8000e003},
		{"lwu t6, 2047(sp)", New().Lwu(xreg(t, 31), xreg(t, 2), off(t, 2047)), 0x7ff16f83},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
