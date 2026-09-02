package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestXorCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"xor a0, a1, a2", New().Xor(xreg(t, 10), xreg(t, 11), xreg(t, 12)), 0x00c5c533},
		{"xor zero, t0, t1", New().Xor(xreg(t, 0), xreg(t, 5), xreg(t, 6)), 0x0062c033},
		{"xor t6, ra, sp", New().Xor(xreg(t, 31), xreg(t, 1), xreg(t, 2)), 0x0020cfb3},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}

	// rd == rs1 in x8-x15 compresses to c.xor (2 bytes).
	b := ctorBytes(t, New().Xor(xreg(t, 10), xreg(t, 10), xreg(t, 11)))
	require.Len(t, b, 2, "xor a0,a0,a1 (c.xor)")
}
