package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddwCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"addw a0, a1, a2", New().Addw(xreg(t, 10), xreg(t, 11), xreg(t, 12)), 0x00c5853b},
		{"addw zero, t0, t1", New().Addw(xreg(t, 0), xreg(t, 5), xreg(t, 6)), 0x0062803b},
		{"addw t6, ra, sp", New().Addw(xreg(t, 31), xreg(t, 1), xreg(t, 2)), 0x00208fbb},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}

	// rd == rs1 in x8-x15 compresses to c.addw (2 bytes).
	b := ctorBytes(t, New().Addw(xreg(t, 10), xreg(t, 10), xreg(t, 11)))
	require.Len(t, b, 2, "addw a0,a0,a1 (c.addw)")
}
