package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSltCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"slt a0, a1, a2", New().Slt(xreg(t, 10), xreg(t, 11), xreg(t, 12)), 0x00c5a533},
		{"slt zero, t0, t1", New().Slt(xreg(t, 0), xreg(t, 5), xreg(t, 6)), 0x0062a033},
		{"slt t6, ra, sp", New().Slt(xreg(t, 31), xreg(t, 1), xreg(t, 2)), 0x0020afb3},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
