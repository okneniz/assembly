package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAuipcCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"auipc a0, 0x12345", New().Auipc(xreg(t, 10), imm20(t, 0x12345)), 0x12345517},
		{"auipc zero, 0", New().Auipc(xreg(t, 0), imm20(t, 0)), 0x00000017},
		{"auipc t0, 0xfffff", New().Auipc(xreg(t, 5), imm20(t, 0xfffff)), 0xfffff297},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
