package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSubwCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"subw a0, a1, a2", New().Subw(xreg(t, 10), xreg(t, 11), xreg(t, 12)), 0x40c5853b},
		{"subw zero, t0, t1", New().Subw(xreg(t, 0), xreg(t, 5), xreg(t, 6)), 0x4062803b},
		{"subw t6, ra, sp", New().Subw(xreg(t, 31), xreg(t, 1), xreg(t, 2)), 0x40208fbb},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}

	// rd == rs1 in x8-x15 compresses to c.subw (2 bytes).
	b := ctorBytes(t, New().Subw(xreg(t, 10), xreg(t, 10), xreg(t, 11)))
	require.Len(t, b, 2, "subw a0,a0,a1 (c.subw)")
}
