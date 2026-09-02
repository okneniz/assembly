package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAmoaddWCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"amoadd.w a0, a2, (a1)", New().AmoaddW(xreg(t, 10), xreg(t, 11), xreg(t, 12)), 0x00c5a52f},
		{"amoadd.w zero, t1, (t0)", New().AmoaddW(xreg(t, 0), xreg(t, 5), xreg(t, 6)), 0x0062a02f},
		{"amoadd.w t6, zero, (t6)", New().AmoaddW(xreg(t, 31), xreg(t, 31), xreg(t, 0)), 0x000fafaf},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
