package loong64

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestALUImmShiftJSONEncodeError - the marshal and write-error paths of
// every immediate ALU and shift instruction (the rest is covered by the
// per-file tests).
func TestALUImmShiftJSONEncodeError(t *testing.T) {
	uimm12, err := New().UImm12(3855)
	require.NoError(t, err)

	imm16, err := New().Imm16(-1)
	require.NoError(t, err)

	uimm5, err := New().UImm5(3)
	require.NoError(t, err)

	uimm6, err := New().UImm6(3)
	require.NoError(t, err)

	family := []struct {
		mnem string
		in   Instr
	}{
		{"addi.d", New().AddiD(lreg(t, 12), lreg(t, 13), imm12v(t, -16))},
		{"slti", New().Slti(lreg(t, 12), lreg(t, 13), imm12v(t, -16))},
		{"sltui", New().Sltui(lreg(t, 12), lreg(t, 13), imm12v(t, -16))},
		{"andi", New().Andi(lreg(t, 12), lreg(t, 13), uimm12)},
		{"ori", New().Ori(lreg(t, 12), lreg(t, 13), uimm12)},
		{"xori", New().Xori(lreg(t, 12), lreg(t, 13), uimm12)},
		{"addu16i.d", New().Addu16iD(lreg(t, 12), lreg(t, 13), imm16)},
		{"slli.w", New().SlliW(lreg(t, 12), lreg(t, 13), uimm5)},
		{"srli.w", New().SrliW(lreg(t, 12), lreg(t, 13), uimm5)},
		{"srai.w", New().SraiW(lreg(t, 12), lreg(t, 13), uimm5)},
		{"rotri.w", New().RotriW(lreg(t, 12), lreg(t, 13), uimm5)},
		{"slli.d", New().SlliD(lreg(t, 12), lreg(t, 13), uimm6)},
		{"srli.d", New().SrliD(lreg(t, 12), lreg(t, 13), uimm6)},
		{"srai.d", New().SraiD(lreg(t, 12), lreg(t, 13), uimm6)},
		{"rotri.d", New().RotriD(lreg(t, 12), lreg(t, 13), uimm6)},
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
