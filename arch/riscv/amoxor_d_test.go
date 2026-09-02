package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAmoxorDCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"amoxor.d a0, a2, (a1)", New().AmoxorD(xreg(t, 10), xreg(t, 11), xreg(t, 12)), 0x20c5b52f},
		{"amoxor.d zero, t1, (t0)", New().AmoxorD(xreg(t, 0), xreg(t, 5), xreg(t, 6)), 0x2062b02f},
		{"amoxor.d t6, zero, (t6)", New().AmoxorD(xreg(t, 31), xreg(t, 31), xreg(t, 0)), 0x200fbfaf},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
