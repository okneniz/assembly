package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCsrrwiCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"csrrwi a0, mstatus, 31", New().Csrrwi(xreg(t, 10), 0x300, 31), 0x300fd573},
		{"csrrwi zero, 0xfff, 0", New().Csrrwi(xreg(t, 0), 0xfff, 0), 0xfff05073},
		{"csrrwi t6, 0, 1", New().Csrrwi(xreg(t, 31), 0, 1), 0x0000dff3},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
