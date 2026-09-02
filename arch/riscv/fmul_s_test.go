package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestFmulSCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"fmul.s fa0, fa1, fa2, rne", New().FmulS(xreg(t, 10), xreg(t, 11), xreg(t, 12), 0), 0x10c58553},
		{"fmul.s ft11, ft0, ft11, dyn", New().FmulS(xreg(t, 31), xreg(t, 0), xreg(t, 31), 7), 0x11f07fd3},
		{"fmul.s ft0, ft11, ft0, rmm", New().FmulS(xreg(t, 0), xreg(t, 31), xreg(t, 0), 4), 0x100fc053},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}

	// The FP registers are named via the FP ABI table (fa0, not a0).
	got, ok := New().FmulS(xreg(t, 10), xreg(t, 11), xreg(t, 12), 0).(FmulS)
	require.True(t, ok)
	require.Equal(
		t,
		"fmul.s fa0, fa1, fa2",
		got.ObjDump(disasm.DefaultViewCtx()),
	)
}
