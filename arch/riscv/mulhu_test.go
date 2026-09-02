package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMulhuCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"mulhu a0, a1, a2", New().Mulhu(xreg(t, 10), xreg(t, 11), xreg(t, 12)), 0x02c5b533},
		{"mulhu zero, t0, t1", New().Mulhu(xreg(t, 0), xreg(t, 5), xreg(t, 6)), 0x0262b033},
		{"mulhu t6, ra, sp", New().Mulhu(xreg(t, 31), xreg(t, 1), xreg(t, 2)), 0x0220bfb3},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
