package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOrCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"or a0, a1, a2", New().Or(xreg(t, 10), xreg(t, 11), xreg(t, 12)), 0x00c5e533},
		{"or zero, t0, t1", New().Or(xreg(t, 0), xreg(t, 5), xreg(t, 6)), 0x0062e033},
		{"or t6, ra, sp", New().Or(xreg(t, 31), xreg(t, 1), xreg(t, 2)), 0x0020efb3},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}

	// rd == rs1 in x8-x15 compresses to c.or (2 bytes).
	b := ctorBytes(t, New().Or(xreg(t, 10), xreg(t, 10), xreg(t, 11)))
	require.Len(t, b, 2, "or a0,a0,a1 (c.or)")
}
