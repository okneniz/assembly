package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAmominuWCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"amominu.w a0, a2, (a1)", New().AmominuW(xreg(t, 10), xreg(t, 11), xreg(t, 12)), 0xc0c5a52f},
		{"amominu.w zero, t1, (t0)", New().AmominuW(xreg(t, 0), xreg(t, 5), xreg(t, 6)), 0xc062a02f},
		{"amominu.w t6, zero, (t6)", New().AmominuW(xreg(t, 31), xreg(t, 31), xreg(t, 0)), 0xc00fafaf},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
