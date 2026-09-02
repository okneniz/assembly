package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSllwCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"sllw a0, a1, a2", New().Sllw(xreg(t, 10), xreg(t, 11), xreg(t, 12)), 0x00c5953b},
		{"sllw zero, t0, t1", New().Sllw(xreg(t, 0), xreg(t, 5), xreg(t, 6)), 0x0062903b},
		{"sllw t6, ra, sp", New().Sllw(xreg(t, 31), xreg(t, 1), xreg(t, 2)), 0x00209fbb},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
