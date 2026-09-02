package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAmoorDCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"amoor.d a0, a2, (a1)", New().AmoorD(xreg(t, 10), xreg(t, 11), xreg(t, 12)), 0x40c5b52f},
		{"amoor.d zero, t1, (t0)", New().AmoorD(xreg(t, 0), xreg(t, 5), xreg(t, 6)), 0x4062b02f},
		{"amoor.d t6, zero, (t6)", New().AmoorD(xreg(t, 31), xreg(t, 31), xreg(t, 0)), 0x400fbfaf},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
