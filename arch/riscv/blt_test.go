package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBltCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"blt a0, a1, 0x1010", New().Blt(xreg(t, 10), xreg(t, 11), 0x1010), 0x00b54863},
		{"blt t0, t1, 0xff8", New().Blt(xreg(t, 5), xreg(t, 6), 0xff8), 0xfe62cce3},
		{"blt a0, a1, 0x1ffe", New().Blt(xreg(t, 10), xreg(t, 11), 0x1ffe), 0x7eb54fe3},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
