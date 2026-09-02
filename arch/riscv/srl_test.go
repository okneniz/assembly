package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSrlCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"srl a0, a1, a2", New().Srl(xreg(t, 10), xreg(t, 11), xreg(t, 12)), 0x00c5d533},
		{"srl zero, t0, t1", New().Srl(xreg(t, 0), xreg(t, 5), xreg(t, 6)), 0x0062d033},
		{"srl t6, ra, sp", New().Srl(xreg(t, 31), xreg(t, 1), xreg(t, 2)), 0x0020dfb3},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
