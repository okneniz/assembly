package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAmoswapWCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"amoswap.w a0, a2, (a1)", New().AmoswapW(xreg(t, 10), xreg(t, 11), xreg(t, 12)), 0x08c5a52f},
		{"amoswap.w zero, t1, (t0)", New().AmoswapW(xreg(t, 0), xreg(t, 5), xreg(t, 6)), 0x0862a02f},
		{"amoswap.w t6, zero, (t6)", New().AmoswapW(xreg(t, 31), xreg(t, 31), xreg(t, 0)), 0x080fafaf},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
