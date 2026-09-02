package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestFdivSCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"fdiv.s fa0, fa1, fa2, rne", New().FdivS(xreg(t, 10), xreg(t, 11), xreg(t, 12), 0), 0x18c58553},
		{"fdiv.s ft11, ft0, ft11, dyn", New().FdivS(xreg(t, 31), xreg(t, 0), xreg(t, 31), 7), 0x19f07fd3},
		{"fdiv.s ft0, ft11, ft0, rmm", New().FdivS(xreg(t, 0), xreg(t, 31), xreg(t, 0), 4), 0x180fc053},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}

	// The FP registers are named via the FP ABI table (fa0, not a0).
	got, ok := New().FdivS(xreg(t, 10), xreg(t, 11), xreg(t, 12), 0).(FdivS)
	require.True(t, ok)
	require.Equal(
		t,
		"fdiv.s fa0, fa1, fa2",
		got.ObjDump(disasm.DefaultViewCtx()),
	)
}
