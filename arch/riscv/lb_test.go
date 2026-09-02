package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLbCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"lb a0, 0(a1)", New().Lb(xreg(t, 10), xreg(t, 11), off(t, 0)), 0x00058503},
		{"lb zero, -2048(ra)", New().Lb(xreg(t, 0), xreg(t, 1), off(t, -2048)), 0x80008003},
		{"lb t6, 2047(sp)", New().Lb(xreg(t, 31), xreg(t, 2), off(t, 2047)), 0x7ff10f83},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
