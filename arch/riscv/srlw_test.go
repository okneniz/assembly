package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSrlwCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"srlw a0, a1, a2", New().Srlw(xreg(t, 10), xreg(t, 11), xreg(t, 12)), 0x00c5d53b},
		{"srlw zero, t0, t1", New().Srlw(xreg(t, 0), xreg(t, 5), xreg(t, 6)), 0x0062d03b},
		{"srlw t6, ra, sp", New().Srlw(xreg(t, 31), xreg(t, 1), xreg(t, 2)), 0x0020dfbb},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
