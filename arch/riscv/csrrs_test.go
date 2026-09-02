package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCsrrsCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"csrrs a0, mstatus, a1", New().Csrrs(xreg(t, 10), 0x300, xreg(t, 11)), 0x3005a573},
		{"csrrs zero, 0xfff, ra", New().Csrrs(xreg(t, 0), 0xfff, xreg(t, 1)), 0xfff0a073},
		{"csrrs t6, 0, zero", New().Csrrs(xreg(t, 31), 0, xreg(t, 0)), 0x00002ff3},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
