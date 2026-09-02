package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCsrrwCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"csrrw a0, mstatus, a1", New().Csrrw(xreg(t, 10), 0x300, xreg(t, 11)), 0x30059573},
		{"csrrw zero, 0xfff, ra", New().Csrrw(xreg(t, 0), 0xfff, xreg(t, 1)), 0xfff09073},
		{"csrrw t6, 0, zero", New().Csrrw(xreg(t, 31), 0, xreg(t, 0)), 0x00001ff3},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
