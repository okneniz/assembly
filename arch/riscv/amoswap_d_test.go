package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAmoswapDCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"amoswap.d a0, a2, (a1)", New().AmoswapD(xreg(t, 10), xreg(t, 11), xreg(t, 12)), 0x08c5b52f},
		{"amoswap.d zero, t1, (t0)", New().AmoswapD(xreg(t, 0), xreg(t, 5), xreg(t, 6)), 0x0862b02f},
		{"amoswap.d t6, zero, (t6)", New().AmoswapD(xreg(t, 31), xreg(t, 31), xreg(t, 0)), 0x080fbfaf},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
