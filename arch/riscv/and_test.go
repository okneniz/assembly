package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAndCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"and a0, a1, a2", New().And(xreg(t, 10), xreg(t, 11), xreg(t, 12)), 0x00c5f533},
		{"and zero, t0, t1", New().And(xreg(t, 0), xreg(t, 5), xreg(t, 6)), 0x0062f033},
		{"and t6, ra, sp", New().And(xreg(t, 31), xreg(t, 1), xreg(t, 2)), 0x0020ffb3},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}

	// rd == rs1 in x8-x15 compresses to c.and (2 bytes).
	b := ctorBytes(t, New().And(xreg(t, 10), xreg(t, 10), xreg(t, 11)))
	require.Len(t, b, 2, "and a0,a0,a1 (c.and)")
}
