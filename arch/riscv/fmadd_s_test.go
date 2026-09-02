package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestFmaddSCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"fmadd.s fa0, fa1, fa2, fa3, rne", New().FmaddS(xreg(t, 10), xreg(t, 11), xreg(t, 12), xreg(t, 13), 0), 0x68c58543},
		{"fmadd.s ft11, ft11, ft11, ft11, dyn", New().FmaddS(xreg(t, 31), xreg(t, 31), xreg(t, 31), xreg(t, 31), 7), 0xf9ffffc3},
		{"fmadd.s ft0, ft1, ft2, ft3, rne", New().FmaddS(xreg(t, 0), xreg(t, 1), xreg(t, 2), xreg(t, 3), 0), 0x18208043},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}

	// The FP registers are named via the FP ABI table (fa0, not a0).
	got, ok := New().FmaddS(xreg(t, 10), xreg(t, 11), xreg(t, 12), xreg(t, 13), 0).(FmaddS)
	require.True(t, ok)
	require.Equal(
		t,
		"fmadd.s fa0, fa1, fa2, fa3",
		got.ObjDump(disasm.DefaultViewCtx()),
	)
}
