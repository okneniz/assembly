package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAmoxorWCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"amoxor.w a0, a2, (a1)", New().AmoxorW(xreg(t, 10), xreg(t, 11), xreg(t, 12)), 0x20c5a52f},
		{"amoxor.w zero, t1, (t0)", New().AmoxorW(xreg(t, 0), xreg(t, 5), xreg(t, 6)), 0x2062a02f},
		{"amoxor.w t6, zero, (t6)", New().AmoxorW(xreg(t, 31), xreg(t, 31), xreg(t, 0)), 0x200fafaf},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
