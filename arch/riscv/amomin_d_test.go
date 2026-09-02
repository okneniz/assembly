package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAmominDCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"amomin.d a0, a2, (a1)", New().AmominD(xreg(t, 10), xreg(t, 11), xreg(t, 12)), 0x80c5b52f},
		{"amomin.d zero, t1, (t0)", New().AmominD(xreg(t, 0), xreg(t, 5), xreg(t, 6)), 0x8062b02f},
		{"amomin.d t6, zero, (t6)", New().AmominD(xreg(t, 31), xreg(t, 31), xreg(t, 0)), 0x800fbfaf},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
