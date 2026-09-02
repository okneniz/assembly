package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAmomaxWCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"amomax.w a0, a2, (a1)", New().AmomaxW(xreg(t, 10), xreg(t, 11), xreg(t, 12)), 0xa0c5a52f},
		{"amomax.w zero, t1, (t0)", New().AmomaxW(xreg(t, 0), xreg(t, 5), xreg(t, 6)), 0xa062a02f},
		{"amomax.w t6, zero, (t6)", New().AmomaxW(xreg(t, 31), xreg(t, 31), xreg(t, 0)), 0xa00fafaf},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
