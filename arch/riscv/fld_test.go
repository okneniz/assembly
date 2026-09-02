package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestFldCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"fld fa0, 8(a1)", New().Fld(xreg(t, 10), xreg(t, 11), off(t, 8)), 0x0085b507},
		{"fld ft0, -2048(zero)", New().Fld(xreg(t, 0), xreg(t, 0), off(t, -2048)), 0x80003007},
		{"fld ft11, 2047(t6)", New().Fld(xreg(t, 31), xreg(t, 31), off(t, 2047)), 0x7fffbf87},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}

	// The FP register is named via the FP ABI table (fa0, not a0).
	got, ok := New().Fld(xreg(t, 10), xreg(t, 11), off(t, 8)).(Fld)
	require.True(t, ok)
	require.Equal(t, "fld fa0, 0x8(a1)",
		got.ObjDump(disasm.DefaultViewCtx()))
}
