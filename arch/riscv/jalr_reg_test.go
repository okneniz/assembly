package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJalrRegCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"jalr t0", New().JalrReg(xreg(t, 5)), 0x000280e7},
		{"jalr a0", New().JalrReg(xreg(t, 10)), 0x000500e7},
		{"jalr ra", New().JalrReg(xreg(t, 1)), 0x000080e7},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
