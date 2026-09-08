package loong64

import (
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
		{"add.w", New().AddW},
		{"add.d", New().AddD},
		{"sub.w", New().SubW},
		{"sub.d", New().SubD},
		{"slt", New().Slt},
		{"sltu", New().Sltu},
		{"maskeqz", New().Maskeqz},
		{"masknez", New().Masknez},
		{"nor", New().Nor},
		{"and", New().And},
		{"or", New().Or},
		{"xor", New().Xor},
		{"orn", New().Orn},
		{"andn", New().Andn},
		{"sll.w", New().SllW},
		{"srl.w", New().SrlW},
		{"sra.w", New().SraW},
		{"sll.d", New().SllD},
		{"srl.d", New().SrlD},
		{"sra.d", New().SraD},
		{"rotr.w", New().RotrW},
		{"rotr.d", New().RotrD},
	}

	for _, f := range family {
		in := f.ctor(lreg(t, 12), lreg(t, 13), lreg(t, 14))

		_, err := in.Encode(errWriter{})
		require.ErrorContains(t, err, "write failed", f.mnem)
	}
}
