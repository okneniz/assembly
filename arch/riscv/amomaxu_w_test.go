package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAmomaxuWCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"amomaxu.w a0, a2, (a1)", New().AmomaxuW(xreg(t, 10), xreg(t, 11), xreg(t, 12)), 0xe0c5a52f},
		{"amomaxu.w zero, t1, (t0)", New().AmomaxuW(xreg(t, 0), xreg(t, 5), xreg(t, 6)), 0xe062a02f},
		{"amomaxu.w t6, zero, (t6)", New().AmomaxuW(xreg(t, 31), xreg(t, 31), xreg(t, 0)), 0xe00fafaf},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
