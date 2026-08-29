package loong64

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestBranchMiscFamilyJSONEncodeError - the marshal and write-error paths
// of every branch/trap instruction of this family (the rest is covered by
// the per-file tests).
func TestBranchMiscFamilyJSONEncodeError(t *testing.T) {
	code0, err := NewCode15(0)
	require.NoError(t, err)

	family := []struct {
		mnem string
		in   Instr
	}{
		{"bne", NewBne(lreg(t, 13), lreg(t, 12), 8)},
		{"blt", NewBlt(lreg(t, 13), lreg(t, 12), 8)},
		{"bge", NewBge(lreg(t, 13), lreg(t, 12), 8)},
		{"bltu", NewBltu(lreg(t, 13), lreg(t, 12), 8)},
		{"bgeu", NewBgeu(lreg(t, 13), lreg(t, 12), 8)},
		{"bnez", NewBnez(lreg(t, 13), 8)},
		{"bl", NewBl(8)},
		{"break", NewBreak(code0)},
		{"syscall", NewSyscall(code0)},
		{"dbar", NewDbar(code0)},
		{"ibar", NewIbar(code0)},
	}

	for _, f := range family {
		b, err := f.in.MarshalJSON()
		require.NoError(t, err, f.mnem)

		var dto map[string]any
		require.NoError(t, json.Unmarshal(b, &dto), f.mnem)
		require.Equal(t, f.mnem, dto["mnemonic"], f.mnem)

		_, err = f.in.Encode(errWriter{}, 0)
		require.ErrorContains(t, err, "write failed", f.mnem)
	}
}
