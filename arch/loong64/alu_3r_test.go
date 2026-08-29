package loong64

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestALU3RFamilyJSONEncodeError - the marshal and write-error paths of
// every 3R ALU instruction (the rest is covered by the per-file tests).
func TestALU3RFamilyJSONEncodeError(t *testing.T) {
	family := []struct {
		mnem string
		ctor func(rd, rj, rk Reg) Instr
	}{
		{"add.w", NewAddW},
		{"add.d", NewAddD},
		{"sub.w", NewSubW},
		{"sub.d", NewSubD},
		{"slt", NewSlt},
		{"sltu", NewSltu},
		{"maskeqz", NewMaskeqz},
		{"masknez", NewMasknez},
		{"nor", NewNor},
		{"and", NewAnd},
		{"or", NewOr},
		{"xor", NewXor},
		{"orn", NewOrn},
		{"andn", NewAndn},
		{"sll.w", NewSllW},
		{"srl.w", NewSrlW},
		{"sra.w", NewSraW},
		{"sll.d", NewSllD},
		{"srl.d", NewSrlD},
		{"sra.d", NewSraD},
		{"rotr.w", NewRotrW},
		{"rotr.d", NewRotrD},
	}

	for _, f := range family {
		in := f.ctor(lreg(t, 12), lreg(t, 13), lreg(t, 14))

		b, err := in.MarshalJSON()
		require.NoError(t, err, f.mnem)

		var dto map[string]any
		require.NoError(t, json.Unmarshal(b, &dto), f.mnem)
		require.Equal(t, f.mnem, dto["mnemonic"], f.mnem)

		_, err = in.Encode(errWriter{}, 0)
		require.ErrorContains(t, err, "write failed", f.mnem)
	}
}
