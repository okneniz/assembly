package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestFswCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"fsw fa2, 8(a1)", New().Fsw(xreg(t, 12), xreg(t, 11), off(t, 8)), 0x00c5a427},
		{"fsw ft0, -2048(t0)", New().Fsw(xreg(t, 0), xreg(t, 5), off(t, -2048)), 0x8002a027},
		{"fsw ft11, 2047(sp)", New().Fsw(xreg(t, 31), xreg(t, 2), off(t, 2047)), 0x7ff12fa7},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}

	// The FP register is named via the FP ABI table (fa2, not a2).
	got, ok := New().Fsw(xreg(t, 12), xreg(t, 11), off(t, 8)).(Fsw)
	require.True(t, ok)
	require.Equal(t, "fsw fa2, 0x8(a1)",
		got.ObjDump(disasm.DefaultViewCtx()))
}
