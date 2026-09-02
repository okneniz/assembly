package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAmominuDCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"amominu.d a0, a2, (a1)", New().AmominuD(xreg(t, 10), xreg(t, 11), xreg(t, 12)), 0xc0c5b52f},
		{"amominu.d zero, t1, (t0)", New().AmominuD(xreg(t, 0), xreg(t, 5), xreg(t, 6)), 0xc062b02f},
		{"amominu.d t6, zero, (t6)", New().AmominuD(xreg(t, 31), xreg(t, 31), xreg(t, 0)), 0xc00fbfaf},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
