package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSltuCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"sltu a0, a1, a2", New().Sltu(xreg(t, 10), xreg(t, 11), xreg(t, 12)), 0x00c5b533},
		{"sltu zero, t0, t1", New().Sltu(xreg(t, 0), xreg(t, 5), xreg(t, 6)), 0x0062b033},
		{"sltu t6, ra, sp", New().Sltu(xreg(t, 31), xreg(t, 1), xreg(t, 2)), 0x0020bfb3},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
