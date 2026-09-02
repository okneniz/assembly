package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSrawCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"sraw a0, a1, a2", New().Sraw(xreg(t, 10), xreg(t, 11), xreg(t, 12)), 0x40c5d53b},
		{"sraw zero, t0, t1", New().Sraw(xreg(t, 0), xreg(t, 5), xreg(t, 6)), 0x4062d03b},
		{"sraw t6, ra, sp", New().Sraw(xreg(t, 31), xreg(t, 1), xreg(t, 2)), 0x4020dfbb},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
