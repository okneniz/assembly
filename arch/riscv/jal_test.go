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
		{"jal ra, 0x1008", New().Jal(xreg(t, 1), 0x1008), 0x008000ef},
		{"jal t0, 0xff8", New().Jal(xreg(t, 5), 0xff8), 0xff9ff2ef},
		{"jal ra, 0x100ffe", New().Jal(xreg(t, 1), 0x100ffe), 0x7ffff0ef},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}

	// rd = zero with a small offset compresses to c.j (2 bytes).
	b := ctorBytes(t, New().Jal(xreg(t, 0), 0x1008))
	require.Len(t, b, 2, "jal zero,0x1008 (c.j)")
}
