package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestFsdCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"fsd fa2, 8(a1)", New().Fsd(xreg(t, 12), xreg(t, 11), off(t, 8)), 0x00c5b427},
		{"fsd ft0, -2048(t0)", New().Fsd(xreg(t, 0), xreg(t, 5), off(t, -2048)), 0x8002b027},
		{"fsd ft11, 2047(sp)", New().Fsd(xreg(t, 31), xreg(t, 2), off(t, 2047)), 0x7ff13fa7},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}

	// The FP register is named via the FP ABI table (fa2, not a2).
	got, ok := New().Fsd(xreg(t, 12), xreg(t, 11), off(t, 8)).(Fsd)
	require.True(t, ok)
	require.Equal(t, "fsd fa2, 0x8(a1)",
		got.ObjDump(disasm.DefaultViewCtx()))
}
