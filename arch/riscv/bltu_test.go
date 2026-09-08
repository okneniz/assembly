package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBltuCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"bltu a0, a1, off 0x10", New().Bltu(xreg(t, 10), xreg(t, 11), 0x10), 0x00b56863},
		{"bltu t0, t1, off -0x8", New().Bltu(xreg(t, 5), xreg(t, 6), -0x8), 0xfe62ece3},
		{"bltu a0, a1, off -0x1000", New().Bltu(xreg(t, 10), xreg(t, 11), -0x1000), 0x80b56063},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
