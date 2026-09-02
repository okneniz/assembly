package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAmoandWCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"amoand.w a0, a2, (a1)", New().AmoandW(xreg(t, 10), xreg(t, 11), xreg(t, 12)), 0x60c5a52f},
		{"amoand.w zero, t1, (t0)", New().AmoandW(xreg(t, 0), xreg(t, 5), xreg(t, 6)), 0x6062a02f},
		{"amoand.w t6, zero, (t6)", New().AmoandW(xreg(t, 31), xreg(t, 31), xreg(t, 0)), 0x600fafaf},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
