package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBgeuCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"bgeu a0, a1, 0x1010", New().Bgeu(xreg(t, 10), xreg(t, 11), 0x1010), 0x00b57863},
		{"bgeu t0, t1, 0xff8", New().Bgeu(xreg(t, 5), xreg(t, 6), 0xff8), 0xfe62fce3},
		{"bgeu a0, a1, 0x1ffe", New().Bgeu(xreg(t, 10), xreg(t, 11), 0x1ffe), 0x7eb57fe3},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}
