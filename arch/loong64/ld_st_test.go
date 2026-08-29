package loong64

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestLdStJSONEncodeError - the marshal and write-error paths of every
// upper-immediate, load and store instruction (the rest is covered by
// the per-file tests).
func TestLdStJSONEncodeError(t *testing.T) {
	imm12v, err := NewImm12(8)
	require.NoError(t, err)

	imm14v, err := NewImm14(8)
	require.NoError(t, err)

	imm20v, err := NewImm20(5)
	require.NoError(t, err)

	uimm5v, err := NewUImm5(5)
	require.NoError(t, err)

	t0, t1, t2 := lreg(t, 12), lreg(t, 13), lreg(t, 14)

	family := []struct {
		mnem string
		in   Instr
	}{
		{"lu32i.d", NewLu32iD(t0, imm20v)},
		{"lu52i.d", NewLu52iD(t0, t1, imm12v)},
		{"pcaddi", NewPcaddi(t0, imm20v)},
		{"pcalau12i", NewPcalau12i(t0, imm20v)},
		{"pcaddu12i", NewPcaddu12i(t0, imm20v)},
		{"pcaddu18i", NewPcaddu18i(t0, imm20v)},
		{"ld.b", NewLdB(t0, t1, imm12v)},
		{"ld.h", NewLdH(t0, t1, imm12v)},
		{"ld.w", NewLdW(t0, t1, imm12v)},
		{"ld.wu", NewLdWu(t0, t1, imm12v)},
		{"ld.d", NewLdD(t0, t1, imm12v)},
		{"st.b", NewStB(t0, t1, imm12v)},
		{"st.h", NewStH(t0, t1, imm12v)},
		{"st.w", NewStW(t0, t1, imm12v)},
		{"st.d", NewStD(t0, t1, imm12v)},
		{"ldptr.w", NewLdptrW(t0, t1, imm14v)},
		{"ldptr.d", NewLdptrD(t0, t1, imm14v)},
		{"stptr.w", NewStptrW(t0, t1, imm14v)},
		{"stptr.d", NewStptrD(t0, t1, imm14v)},
		{"ldx.b", NewLdxB(t0, t1, t2)},
		{"ldx.h", NewLdxH(t0, t1, t2)},
		{"ldx.w", NewLdxW(t0, t1, t2)},
		{"ldx.wu", NewLdxWu(t0, t1, t2)},
		{"ldx.d", NewLdxD(t0, t1, t2)},
		{"stx.b", NewStxB(t0, t1, t2)},
		{"stx.h", NewStxH(t0, t1, t2)},
		{"stx.w", NewStxW(t0, t1, t2)},
		{"stx.d", NewStxD(t0, t1, t2)},
		{"preld", NewPreld(uimm5v, t1, imm12v)},
		{"preldx", NewPreldx(uimm5v, t1, t2)},
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
