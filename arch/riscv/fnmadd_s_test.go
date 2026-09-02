package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestFnmaddSCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"fnmadd.s fa0, fa1, fa2, fa3, rne", New().FnmaddS(xreg(t, 10), xreg(t, 11), xreg(t, 12), xreg(t, 13), 0), 0x68c5854f},
		{"fnmadd.s ft11, ft11, ft11, ft11, dyn", New().FnmaddS(xreg(t, 31), xreg(t, 31), xreg(t, 31), xreg(t, 31), 7), 0xf9ffffcf},
		{"fnmadd.s ft0, ft1, ft2, ft3, rne", New().FnmaddS(xreg(t, 0), xreg(t, 1), xreg(t, 2), xreg(t, 3), 0), 0x1820804f},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}

	// The FP registers are named via the FP ABI table (fa0, not a0).
	got, ok := New().FnmaddS(xreg(t, 10), xreg(t, 11), xreg(t, 12), xreg(t, 13), 0).(FnmaddS)
	require.True(t, ok)
	require.Equal(
		t,
		"fnmadd.s fa0, fa1, fa2, fa3",
		got.ObjDump(disasm.DefaultViewCtx()),
	)
}
