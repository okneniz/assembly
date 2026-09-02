package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAmominWCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"amomin.w a0, a2, (a1)", New().AmominW(xreg(t, 10), xreg(t, 11), xreg(t, 12)), 0x80c5a52f},
		{"amomin.w zero, t1, (t0)", New().AmominW(xreg(t, 0), xreg(t, 5), xreg(t, 6)), 0x8062a02f},
		{"amomin.w t6, zero, (t6)", New().AmominW(xreg(t, 31), xreg(t, 31), xreg(t, 0)), 0x800fafaf},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
