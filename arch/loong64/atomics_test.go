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
	off, err := NewImm14(8)
	require.NoError(t, err)

	family := []struct {
		mnem string
		in   Instr
	}{
		{"ll.w", NewLlW(lreg(t, 12), lreg(t, 13), off)},
		{"ll.d", NewLlD(lreg(t, 12), lreg(t, 13), off)},
		{"sc.w", NewScW(lreg(t, 12), lreg(t, 13), off)},
		{"sc.d", NewScD(lreg(t, 12), lreg(t, 13), off)},
		{"llacq.w", NewLlacqW(lreg(t, 12), lreg(t, 13))},
		{"llacq.d", NewLlacqD(lreg(t, 12), lreg(t, 13))},
		{"screl.w", NewScrelW(lreg(t, 12), lreg(t, 13))},
		{"screl.d", NewScrelD(lreg(t, 12), lreg(t, 13))},
		{"sc.q", NewScQ(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amcas.b", NewAmcasB(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amcas.h", NewAmcasH(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amcas.w", NewAmcasW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amcas.d", NewAmcasD(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amcas_db.b", NewAmcasDbB(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amcas_db.h", NewAmcasDbH(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amcas_db.w", NewAmcasDbW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amcas_db.d", NewAmcasDbD(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amswap.b", NewAmswapB(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amswap.h", NewAmswapH(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amswap.w", NewAmswapW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amswap.d", NewAmswapD(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amswap_db.b", NewAmswapDbB(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amswap_db.h", NewAmswapDbH(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amswap_db.w", NewAmswapDbW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amswap_db.d", NewAmswapDbD(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amadd.b", NewAmaddB(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amadd.h", NewAmaddH(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amadd.w", NewAmaddW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amadd.d", NewAmaddD(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amadd_db.b", NewAmaddDbB(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amadd_db.h", NewAmaddDbH(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amadd_db.w", NewAmaddDbW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amadd_db.d", NewAmaddDbD(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amand.w", NewAmandW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amand.d", NewAmandD(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amand_db.w", NewAmandDbW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amand_db.d", NewAmandDbD(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amor.w", NewAmorW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amor.d", NewAmorD(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amor_db.w", NewAmorDbW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amor_db.d", NewAmorDbD(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amxor.w", NewAmxorW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amxor.d", NewAmxorD(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amxor_db.w", NewAmxorDbW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"amxor_db.d", NewAmxorDbD(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"ammax.w", NewAmmaxW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"ammax.d", NewAmmaxD(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"ammax.wu", NewAmmaxWu(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"ammax.du", NewAmmaxDu(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"ammax_db.w", NewAmmaxDbW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"ammax_db.d", NewAmmaxDbD(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"ammax_db.wu", NewAmmaxDbWu(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"ammax_db.du", NewAmmaxDbDu(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"ammin.w", NewAmminW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"ammin.d", NewAmminD(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"ammin.wu", NewAmminWu(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"ammin.du", NewAmminDu(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"ammin_db.w", NewAmminDbW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"ammin_db.d", NewAmminDbD(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"ammin_db.wu", NewAmminDbWu(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"ammin_db.du", NewAmminDbDu(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
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
