package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSraiwCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"sraiw a0, a1, 0", New().Sraiw(xreg(t, 10), xreg(t, 11), imm12(t, 0)), 0x4005d51b},
		{"sraiw a0, a1, 1", New().Sraiw(xreg(t, 10), xreg(t, 11), imm12(t, 1)), 0x4015d51b},
		{"sraiw a0, a1, 31", New().Sraiw(xreg(t, 10), xreg(t, 11), imm12(t, 31)), 0x41f5d51b},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
