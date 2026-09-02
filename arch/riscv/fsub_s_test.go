package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestFsubSCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"fsub.s fa0, fa1, fa2, rne", New().FsubS(xreg(t, 10), xreg(t, 11), xreg(t, 12), 0), 0x08c58553},
		{"fsub.s ft11, ft0, ft11, dyn", New().FsubS(xreg(t, 31), xreg(t, 0), xreg(t, 31), 7), 0x09f07fd3},
		{"fsub.s ft0, ft11, ft0, rmm", New().FsubS(xreg(t, 0), xreg(t, 31), xreg(t, 0), 4), 0x080fc053},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}

	// The FP registers are named via the FP ABI table (fa0, not a0).
	got, ok := New().FsubS(xreg(t, 10), xreg(t, 11), xreg(t, 12), 0).(FsubS)
	require.True(t, ok)
	require.Equal(
		t,
		"fsub.s fa0, fa1, fa2",
		got.ObjDump(disasm.DefaultViewCtx()),
	)
}
