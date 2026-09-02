package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCsrrcCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"csrrc a0, mstatus, a1", New().Csrrc(xreg(t, 10), 0x300, xreg(t, 11)), 0x3005b573},
		{"csrrc zero, 0xfff, ra", New().Csrrc(xreg(t, 0), 0xfff, xreg(t, 1)), 0xfff0b073},
		{"csrrc t6, 0, zero", New().Csrrc(xreg(t, 31), 0, xreg(t, 0)), 0x00003ff3},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
