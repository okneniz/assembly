package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAmoandDCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"amoand.d a0, a2, (a1)", New().AmoandD(xreg(t, 10), xreg(t, 11), xreg(t, 12)), 0x60c5b52f},
		{"amoand.d zero, t1, (t0)", New().AmoandD(xreg(t, 0), xreg(t, 5), xreg(t, 6)), 0x6062b02f},
		{"amoand.d t6, zero, (t6)", New().AmoandD(xreg(t, 31), xreg(t, 31), xreg(t, 0)), 0x600fbfaf},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
