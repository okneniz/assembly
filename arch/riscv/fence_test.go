package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFenceCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"fence 0x0", New().Fence(0x0), 0x0000000f},
		{"fence 0x3", New().Fence(0x3), 0x0030000f},
		{"fence 0xf", New().Fence(0xf), 0x00f0000f},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
