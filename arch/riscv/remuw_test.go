package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRemuwCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"remuw a0, a1, a2", New().Remuw(xreg(t, 10), xreg(t, 11), xreg(t, 12)), 0x02c5f53b},
		{"remuw zero, t0, t1", New().Remuw(xreg(t, 0), xreg(t, 5), xreg(t, 6)), 0x0262f03b},
		{"remuw t6, ra, sp", New().Remuw(xreg(t, 31), xreg(t, 1), xreg(t, 2)), 0x0220ffbb},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
