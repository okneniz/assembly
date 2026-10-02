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
		{"fence 0x0", New().Fence(0x0), 0x0ff0000f}, // llvm-mc: iorw|iorw
		{"fence 0x3", New().Fence(0x3), 0x3ff0000f},
		{"fence 0xf", New().Fence(0xf), 0xfff0000f},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
