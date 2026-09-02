package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAmoaddDCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"amoadd.d a0, a2, (a1)", New().AmoaddD(xreg(t, 10), xreg(t, 11), xreg(t, 12)), 0x00c5b52f},
		{"amoadd.d zero, t1, (t0)", New().AmoaddD(xreg(t, 0), xreg(t, 5), xreg(t, 6)), 0x0062b02f},
		{"amoadd.d t6, zero, (t6)", New().AmoaddD(xreg(t, 31), xreg(t, 31), xreg(t, 0)), 0x000fbfaf},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
