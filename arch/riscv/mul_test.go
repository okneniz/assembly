package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMulCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"mul a0, a1, a2", New().Mul(xreg(t, 10), xreg(t, 11), xreg(t, 12)), 0x02c58533},
		{"mul zero, t0, t1", New().Mul(xreg(t, 0), xreg(t, 5), xreg(t, 6)), 0x02628033},
		{"mul t6, ra, sp", New().Mul(xreg(t, 31), xreg(t, 1), xreg(t, 2)), 0x02208fb3},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
