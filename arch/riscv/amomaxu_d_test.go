package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAmomaxuDCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"amomaxu.d a0, a2, (a1)", New().AmomaxuD(xreg(t, 10), xreg(t, 11), xreg(t, 12)), 0xe0c5b52f},
		{"amomaxu.d zero, t1, (t0)", New().AmomaxuD(xreg(t, 0), xreg(t, 5), xreg(t, 6)), 0xe062b02f},
		{"amomaxu.d t6, zero, (t6)", New().AmomaxuD(xreg(t, 31), xreg(t, 31), xreg(t, 0)), 0xe00fbfaf},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
