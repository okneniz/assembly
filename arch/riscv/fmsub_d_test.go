package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestFmsubDCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"fmsub.d fa0, fa1, fa2, fa3, rne", New().FmsubD(xreg(t, 10), xreg(t, 11), xreg(t, 12), xreg(t, 13), 0), 0x6ac58547},
		{"fmsub.d ft11, ft11, ft11, ft11, dyn", New().FmsubD(xreg(t, 31), xreg(t, 31), xreg(t, 31), xreg(t, 31), 7), 0xfbffffc7},
		{"fmsub.d ft0, ft1, ft2, ft3, rne", New().FmsubD(xreg(t, 0), xreg(t, 1), xreg(t, 2), xreg(t, 3), 0), 0x1a208047},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}

	// The FP registers are named via the FP ABI table (fa0, not a0).
	got, ok := New().FmsubD(xreg(t, 10), xreg(t, 11), xreg(t, 12), xreg(t, 13), 0).(FmsubD)
	require.True(t, ok)
	require.Equal(
		t,
		"fmsub.d fa0, fa1, fa2, fa3",
		got.ObjDump(disasm.DefaultViewCtx()),
	)
}
