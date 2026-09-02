package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestXoriCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"xori a0, a1, 2047", New().Xori(xreg(t, 10), xreg(t, 11), imm12(t, 2047)), 0x7ff5c513},
		{"xori a0, a1, -2048", New().Xori(xreg(t, 10), xreg(t, 11), imm12(t, -2048)), 0x8005c513},
		{"xori a0, a1, -1", New().Xori(xreg(t, 10), xreg(t, 11), imm12(t, -1)), 0xfff5c513},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
