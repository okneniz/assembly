package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJalCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"jal ra, off 0x8", New().Jal(xreg(t, 1), 0x8), 0x008000ef},
		{"jal t0, off -0x8", New().Jal(xreg(t, 5), -0x8), 0xff9ff2ef},
		{"jal ra, off 0xffffe", New().Jal(xreg(t, 1), 0xffffe), 0x7ffff0ef},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}

	// rd = zero with a small offset compresses to c.j (2 bytes).
	b := ctorBytes(t, New().Jal(xreg(t, 0), 0x8))
	require.Len(t, b, 2, "jal zero,0x1008 (c.j)")
}
