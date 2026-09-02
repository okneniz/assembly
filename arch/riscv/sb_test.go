package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSbCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"sb a0, 0(a1)", New().Sb(xreg(t, 10), xreg(t, 11), off(t, 0)), 0x00a58023},
		{"sb ra, -2048(sp)", New().Sb(xreg(t, 1), xreg(t, 2), off(t, -2048)), 0x80110023},
		{"sb t6, 2047(sp)", New().Sb(xreg(t, 31), xreg(t, 2), off(t, 2047)), 0x7ff10fa3},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
