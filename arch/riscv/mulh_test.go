package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMulhCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"mulh a0, a1, a2", New().Mulh(xreg(t, 10), xreg(t, 11), xreg(t, 12)), 0x02c59533},
		{"mulh zero, t0, t1", New().Mulh(xreg(t, 0), xreg(t, 5), xreg(t, 6)), 0x02629033},
		{"mulh t6, ra, sp", New().Mulh(xreg(t, 31), xreg(t, 1), xreg(t, 2)), 0x02209fb3},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
