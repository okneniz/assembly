package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestBranchMiscFamilyJSONEncodeError - the marshal and write-error paths
// of every branch/trap instruction of this family (the rest is covered by
// the per-file tests).
func TestBranchMiscFamilyJSONEncodeError(t *testing.T) {
	code0, err := New().Code15(0)
	require.NoError(t, err)

	family := []struct {
		mnem string
		in   Instr
	}{
		{"bne", New().Bne(lreg(t, 13), lreg(t, 12), 8)},
		{"blt", New().Blt(lreg(t, 13), lreg(t, 12), 8)},
		{"bge", New().Bge(lreg(t, 13), lreg(t, 12), 8)},
		{"bltu", New().Bltu(lreg(t, 13), lreg(t, 12), 8)},
		{"bgeu", New().Bgeu(lreg(t, 13), lreg(t, 12), 8)},
		{"bnez", New().Bnez(lreg(t, 13), 8)},
		{"bl", New().Bl(8)},
		{"break", New().Break(code0)},
		{"syscall", New().Syscall(code0)},
		{"dbar", New().Dbar(code0)},
		{"ibar", New().Ibar(code0)},
	}

	for _, f := range family {
		_, err := f.in.Encode(errWriter{}, 0)
		require.ErrorContains(t, err, "write failed", f.mnem)
	}
}
