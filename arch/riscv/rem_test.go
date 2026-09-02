package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRemCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"rem a0, a1, a2", New().Rem(xreg(t, 10), xreg(t, 11), xreg(t, 12)), 0x02c5e533},
		{"rem zero, t0, t1", New().Rem(xreg(t, 0), xreg(t, 5), xreg(t, 6)), 0x0262e033},
		{"rem t6, ra, sp", New().Rem(xreg(t, 31), xreg(t, 1), xreg(t, 2)), 0x0220efb3},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
