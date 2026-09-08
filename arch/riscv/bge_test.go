package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBgeCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"bge a0, a1, off 0x10", New().Bge(xreg(t, 10), xreg(t, 11), 0x10), 0x00b55863},
		{"bge t0, t1, off -0x8", New().Bge(xreg(t, 5), xreg(t, 6), -0x8), 0xfe62dce3},
		{"bge a0, a1, off -0x1000", New().Bge(xreg(t, 10), xreg(t, 11), -0x1000), 0x80b55063},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
