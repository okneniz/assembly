package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAmomaxDCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"amomax.d a0, a2, (a1)", New().AmomaxD(xreg(t, 10), xreg(t, 11), xreg(t, 12)), 0xa0c5b52f},
		{"amomax.d zero, t1, (t0)", New().AmomaxD(xreg(t, 0), xreg(t, 5), xreg(t, 6)), 0xa062b02f},
		{"amomax.d t6, zero, (t6)", New().AmomaxD(xreg(t, 31), xreg(t, 31), xreg(t, 0)), 0xa00fbfaf},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
