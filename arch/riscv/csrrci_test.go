package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCsrrciCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"csrrci a0, mstatus, 31", New().Csrrci(xreg(t, 10), 0x300, 31), 0x300ff573},
		{"csrrci zero, 0xfff, 0", New().Csrrci(xreg(t, 0), 0xfff, 0), 0xfff07073},
		{"csrrci t6, 0, 1", New().Csrrci(xreg(t, 31), 0, 1), 0x0000fff3},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
