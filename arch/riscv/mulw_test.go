package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMulwCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"mulw a0, a1, a2", New().Mulw(xreg(t, 10), xreg(t, 11), xreg(t, 12)), 0x02c5853b},
		{"mulw zero, t0, t1", New().Mulw(xreg(t, 0), xreg(t, 5), xreg(t, 6)), 0x0262803b},
		{"mulw t6, ra, sp", New().Mulw(xreg(t, 31), xreg(t, 1), xreg(t, 2)), 0x02208fbb},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
