package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSrliwCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"srliw a0, a1, 0", New().Srliw(xreg(t, 10), xreg(t, 11), imm12(t, 0)), 0x0005d51b},
		{"srliw a0, a1, 1", New().Srliw(xreg(t, 10), xreg(t, 11), imm12(t, 1)), 0x0015d51b},
		{"srliw a0, a1, 31", New().Srliw(xreg(t, 10), xreg(t, 11), imm12(t, 31)), 0x01f5d51b},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
