package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestFdivDCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"fdiv.d fa0, fa1, fa2, rne", New().FdivD(xreg(t, 10), xreg(t, 11), xreg(t, 12), 0), 0x1ac58553},
		{"fdiv.d ft11, ft0, ft11, dyn", New().FdivD(xreg(t, 31), xreg(t, 0), xreg(t, 31), 7), 0x1bf07fd3},
		{"fdiv.d ft0, ft11, ft0, rmm", New().FdivD(xreg(t, 0), xreg(t, 31), xreg(t, 0), 4), 0x1a0fc053},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}

	// The FP registers are named via the FP ABI table (fa0, not a0).
	got, ok := New().FdivD(xreg(t, 10), xreg(t, 11), xreg(t, 12), 0).(FdivD)
	require.True(t, ok)
	require.Equal(
		t,
		"fdiv.d fa0, fa1, fa2",
		got.ObjDump(disasm.DefaultViewCtx()),
	)
}
