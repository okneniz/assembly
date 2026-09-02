package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestFaddSCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"fadd.s fa0, fa1, fa2, rne", New().FaddS(xreg(t, 10), xreg(t, 11), xreg(t, 12), 0), 0x00c58553},
		{"fadd.s ft11, ft0, ft11, dyn", New().FaddS(xreg(t, 31), xreg(t, 0), xreg(t, 31), 7), 0x01f07fd3},
		{"fadd.s ft0, ft11, ft0, rmm", New().FaddS(xreg(t, 0), xreg(t, 31), xreg(t, 0), 4), 0x000fc053},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}

	// The FP registers are named via the FP ABI table (fa0, not a0).
	got, ok := New().FaddS(xreg(t, 10), xreg(t, 11), xreg(t, 12), 0).(FaddS)
	require.True(t, ok)
	require.Equal(
		t,
		"fadd.s fa0, fa1, fa2",
		got.ObjDump(disasm.DefaultViewCtx()),
	)
}
