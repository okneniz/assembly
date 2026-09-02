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
	imm12v, err := New().Imm12(8)
	require.NoError(t, err)

	imm14v, err := New().Imm14(8)
	require.NoError(t, err)

	imm20v, err := New().Imm20(5)
	require.NoError(t, err)

	uimm5v, err := New().UImm5(5)
	require.NoError(t, err)

	t0, t1, t2 := lreg(t, 12), lreg(t, 13), lreg(t, 14)

	family := []struct {
		mnem string
		in   Instr
	}{
		{"lu32i.d", New().Lu32iD(t0, imm20v)},
		{"lu52i.d", New().Lu52iD(t0, t1, imm12v)},
		{"pcaddi", New().Pcaddi(t0, imm20v)},
		{"pcalau12i", New().Pcalau12i(t0, imm20v)},
		{"pcaddu12i", New().Pcaddu12i(t0, imm20v)},
		{"pcaddu18i", New().Pcaddu18i(t0, imm20v)},
		{"ld.b", New().LdB(t0, t1, imm12v)},
		{"ld.h", New().LdH(t0, t1, imm12v)},
		{"ld.w", New().LdW(t0, t1, imm12v)},
		{"ld.wu", New().LdWu(t0, t1, imm12v)},
		{"ld.d", New().LdD(t0, t1, imm12v)},
		{"st.b", New().StB(t0, t1, imm12v)},
		{"st.h", New().StH(t0, t1, imm12v)},
		{"st.w", New().StW(t0, t1, imm12v)},
		{"st.d", New().StD(t0, t1, imm12v)},
		{"ldptr.w", New().LdptrW(t0, t1, imm14v)},
		{"ldptr.d", New().LdptrD(t0, t1, imm14v)},
		{"stptr.w", New().StptrW(t0, t1, imm14v)},
		{"stptr.d", New().StptrD(t0, t1, imm14v)},
		{"ldx.b", New().LdxB(t0, t1, t2)},
		{"ldx.h", New().LdxH(t0, t1, t2)},
		{"ldx.w", New().LdxW(t0, t1, t2)},
		{"ldx.wu", New().LdxWu(t0, t1, t2)},
		{"ldx.d", New().LdxD(t0, t1, t2)},
		{"stx.b", New().StxB(t0, t1, t2)},
		{"stx.h", New().StxH(t0, t1, t2)},
		{"stx.w", New().StxW(t0, t1, t2)},
		{"stx.d", New().StxD(t0, t1, t2)},
		{"preld", New().Preld(uimm5v, t1, imm12v)},
		{"preldx", New().Preldx(uimm5v, t1, t2)},
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
