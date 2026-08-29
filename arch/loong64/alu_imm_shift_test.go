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
	uimm12, err := NewUImm12(3855)
	require.NoError(t, err)

	imm16, err := NewImm16(-1)
	require.NoError(t, err)

	uimm5, err := NewUImm5(3)
	require.NoError(t, err)

	uimm6, err := NewUImm6(3)
	require.NoError(t, err)

	family := []struct {
		mnem string
		in   Instr
	}{
		{"addi.d", NewAddiD(lreg(t, 12), lreg(t, 13), imm12v(t, -16))},
		{"slti", NewSlti(lreg(t, 12), lreg(t, 13), imm12v(t, -16))},
		{"sltui", NewSltui(lreg(t, 12), lreg(t, 13), imm12v(t, -16))},
		{"andi", NewAndi(lreg(t, 12), lreg(t, 13), uimm12)},
		{"ori", NewOri(lreg(t, 12), lreg(t, 13), uimm12)},
		{"xori", NewXori(lreg(t, 12), lreg(t, 13), uimm12)},
		{"addu16i.d", NewAddu16iD(lreg(t, 12), lreg(t, 13), imm16)},
		{"slli.w", NewSlliW(lreg(t, 12), lreg(t, 13), uimm5)},
		{"srli.w", NewSrliW(lreg(t, 12), lreg(t, 13), uimm5)},
		{"srai.w", NewSraiW(lreg(t, 12), lreg(t, 13), uimm5)},
		{"rotri.w", NewRotriW(lreg(t, 12), lreg(t, 13), uimm5)},
		{"slli.d", NewSlliD(lreg(t, 12), lreg(t, 13), uimm6)},
		{"srli.d", NewSrliD(lreg(t, 12), lreg(t, 13), uimm6)},
		{"srai.d", NewSraiD(lreg(t, 12), lreg(t, 13), uimm6)},
		{"rotri.d", NewRotriD(lreg(t, 12), lreg(t, 13), uimm6)},
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
