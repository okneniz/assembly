package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCsrrsiCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"csrrsi a0, mstatus, 31", New().Csrrsi(xreg(t, 10), 0x300, 31), 0x300fe573},
		{"csrrsi zero, 0xfff, 0", New().Csrrsi(xreg(t, 0), 0xfff, 0), 0xfff06073},
		{"csrrsi t6, 0, 1", New().Csrrsi(xreg(t, 31), 0, 1), 0x0000eff3},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
