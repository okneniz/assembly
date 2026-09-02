package loong64

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestAtomicsJSONEncodeError - the marshal and write-error paths of every
// atomics instruction (ll/sc, acquire/release and the AMO family; the rest
// is covered by the per-file tests).
func TestAtomicsJSONEncodeError(t *testing.T) {
	off, err := New().Imm14(8)
	require.NoError(t, err)

	family := []struct {
		mnem string
		in   Instr
	}{
		{"ll.w", New().LlW(lreg(t, 12), lreg(t, 13), off)},
		{"ll.d", New().LlD(lreg(t, 12), lreg(t, 13), off)},
		{"sc.w", New().ScW(lreg(t, 12), lreg(t, 13), off)},
		{"sc.d", New().ScD(lreg(t, 12), lreg(t, 13), off)},
		{"llacq.w", New().LlacqW(lreg(t, 12), lreg(t, 13))},
		{"llacq.d", New().LlacqD(lreg(t, 12), lreg(t, 13))},
		{"screl.w", New().ScrelW(lreg(t, 12), lreg(t, 13))},
		{"screl.d", New().ScrelD(lreg(t, 12), lreg(t, 13))},
		{"sc.q", New().ScQ(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amcas.b", New().AmcasB(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amcas.h", New().AmcasH(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amcas.w", New().AmcasW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amcas.d", New().AmcasD(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amcas_db.b", New().AmcasDbB(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amcas_db.h", New().AmcasDbH(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amcas_db.w", New().AmcasDbW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amcas_db.d", New().AmcasDbD(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amswap.b", New().AmswapB(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amswap.h", New().AmswapH(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amswap.w", New().AmswapW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amswap.d", New().AmswapD(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amswap_db.b", New().AmswapDbB(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amswap_db.h", New().AmswapDbH(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amswap_db.w", New().AmswapDbW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amswap_db.d", New().AmswapDbD(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amadd.b", New().AmaddB(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amadd.h", New().AmaddH(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amadd.w", New().AmaddW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amadd.d", New().AmaddD(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amadd_db.b", New().AmaddDbB(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amadd_db.h", New().AmaddDbH(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amadd_db.w", New().AmaddDbW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amadd_db.d", New().AmaddDbD(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amand.w", New().AmandW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amand.d", New().AmandD(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amand_db.w", New().AmandDbW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amand_db.d", New().AmandDbD(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amor.w", New().AmorW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amor.d", New().AmorD(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amor_db.w", New().AmorDbW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amor_db.d", New().AmorDbD(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amxor.w", New().AmxorW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amxor.d", New().AmxorD(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amxor_db.w", New().AmxorDbW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amxor_db.d", New().AmxorDbD(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"ammax.w", New().AmmaxW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"ammax.d", New().AmmaxD(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"ammax.wu", New().AmmaxWu(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"ammax.du", New().AmmaxDu(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"ammax_db.w", New().AmmaxDbW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"ammax_db.d", New().AmmaxDbD(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"ammax_db.wu", New().AmmaxDbWu(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"ammax_db.du", New().AmmaxDbDu(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"ammin.w", New().AmminW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"ammin.d", New().AmminD(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"ammin.wu", New().AmminWu(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"ammin.du", New().AmminDu(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"ammin_db.w", New().AmminDbW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"ammin_db.d", New().AmminDbD(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"ammin_db.wu", New().AmminDbWu(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"ammin_db.du", New().AmminDbDu(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
	}

	for _, f := range family {
		b, err := f.in.MarshalJSON()
		require.NoError(t, err, f.mnem)

		var dto map[string]any
		require.NoError(t, json.Unmarshal(b, &dto), f.mnem)
		require.Equal(t, f.mnem, dto["mnemonic"], f.mnem)
		require.Equal(t, "LA64", dto["group"], f.mnem)

		_, err = f.in.Encode(errWriter{}, 0)
		require.ErrorContains(t, err, "write failed", f.mnem)
	}
}
