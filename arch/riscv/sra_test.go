package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSraCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"sra a0, a1, a2", New().Sra(xreg(t, 10), xreg(t, 11), xreg(t, 12)), 0x40c5d533},
		{"sra zero, t0, t1", New().Sra(xreg(t, 0), xreg(t, 5), xreg(t, 6)), 0x4062d033},
		{"sra t6, ra, sp", New().Sra(xreg(t, 31), xreg(t, 1), xreg(t, 2)), 0x4020dfb3},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
