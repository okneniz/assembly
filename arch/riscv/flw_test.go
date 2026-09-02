package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestFlwCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"flw fa0, 8(a1)", New().Flw(xreg(t, 10), xreg(t, 11), off(t, 8)), 0x0085a507},
		{"flw ft0, -2048(zero)", New().Flw(xreg(t, 0), xreg(t, 0), off(t, -2048)), 0x80002007},
		{"flw ft11, 2047(t6)", New().Flw(xreg(t, 31), xreg(t, 31), off(t, 2047)), 0x7fffaf87},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}

	// The FP register is named via the FP ABI table (fa0, not a0).
	got, ok := New().Flw(xreg(t, 10), xreg(t, 11), off(t, 8)).(Flw)
	require.True(t, ok)
	require.Equal(t, "flw fa0, 0x8(a1)",
		got.ObjDump(disasm.DefaultViewCtx()))
}
