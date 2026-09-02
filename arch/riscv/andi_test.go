package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAndiCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"andi a0, a1, 2047", New().Andi(xreg(t, 10), xreg(t, 11), imm12(t, 2047)), 0x7ff5f513},
		{"andi a0, a1, -2048", New().Andi(xreg(t, 10), xreg(t, 11), imm12(t, -2048)), 0x8005f513},
		{"andi a0, a1, 0xff", New().Andi(xreg(t, 10), xreg(t, 11), imm12(t, 0xff)), 0x0ff5f513},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}

	// rd == rs1 in x8-x15 with a 6-bit imm compresses to c.andi (2 bytes).
	b := ctorBytes(t, New().Andi(xreg(t, 10), xreg(t, 10), imm12(t, 1)))
	require.Len(t, b, 2, "andi a0,a0,1 (c.andi)")
}
