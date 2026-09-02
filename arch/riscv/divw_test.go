package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDivwCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"divw a0, a1, a2", New().Divw(xreg(t, 10), xreg(t, 11), xreg(t, 12)), 0x02c5c53b},
		{"divw zero, t0, t1", New().Divw(xreg(t, 0), xreg(t, 5), xreg(t, 6)), 0x0262c03b},
		{"divw t6, ra, sp", New().Divw(xreg(t, 31), xreg(t, 1), xreg(t, 2)), 0x0220cfbb},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
