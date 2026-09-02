package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"sh a0, 0(a1)", New().Sh(xreg(t, 10), xreg(t, 11), off(t, 0)), 0x00a59023},
		{"sh ra, -2048(sp)", New().Sh(xreg(t, 1), xreg(t, 2), off(t, -2048)), 0x80111023},
		{"sh t6, 2047(sp)", New().Sh(xreg(t, 31), xreg(t, 2), off(t, 2047)), 0x7ff11fa3},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
