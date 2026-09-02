package loong64

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestMul2RJSONEncodeError - the marshal and write-error paths of every
// mul/div/mod (3R) and ext/rdtime/cpucfg/bit-ops (2R) instruction (the
// rest is covered by the per-file tests).
func TestMul2RJSONEncodeError(t *testing.T) {
	family3r := []struct {
		mnem string
		ctor func(rd, rj, rk Reg) Instr
	}{
		{"mul.w", New().MulW},
		{"mulh.w", New().MulhW},
		{"mulh.wu", New().MulhWu},
		{"mul.d", New().MulD},
		{"mulh.d", New().MulhD},
		{"mulh.du", New().MulhDu},
		{"mulw.d.w", New().MulwDW},
		{"mulw.d.wu", New().MulwDWu},
		{"div.w", New().DivW},
		{"mod.w", New().ModW},
		{"div.wu", New().DivWu},
		{"mod.wu", New().ModWu},
		{"div.d", New().DivD},
		{"mod.d", New().ModD},
		{"div.du", New().DivDu},
		{"mod.du", New().ModDu},
	}

	family2r := []struct {
		mnem string
		ctor func(rd, rj Reg) Instr
	}{
		{"ext.w.b", New().ExtWB},
		{"ext.w.h", New().ExtWH},
		{"rdtimel.w", New().RdtimelW},
		{"rdtimeh.w", New().RdtimehW},
		{"rdtime.d", New().RdtimeD},
		{"cpucfg", New().Cpucfg},
		{"clo.w", New().CloW},
		{"clz.w", New().ClzW},
		{"cto.w", New().CtoW},
		{"ctz.w", New().CtzW},
		{"clo.d", New().CloD},
		{"clz.d", New().ClzD},
		{"cto.d", New().CtoD},
		{"ctz.d", New().CtzD},
		{"revb.2h", New().Revb2H},
		{"revb.4h", New().Revb4H},
		{"revb.2w", New().Revb2W},
		{"revb.d", New().RevbD},
		{"revh.2w", New().Revh2W},
		{"revh.d", New().RevhD},
		{"bitrev.4b", New().Revbit4B},
		{"bitrev.w", New().RevbitW},
		{"bitrev.8b", New().Revbit8B},
		{"bitrev.d", New().RevbitD},
	}

	check := func(mnem string, in Instr) {
		b, err := in.MarshalJSON()
		require.NoError(t, err, mnem)

		var dto map[string]any
		require.NoError(t, json.Unmarshal(b, &dto), mnem)
		require.Equal(t, mnem, dto["mnemonic"], mnem)

		_, err = in.Encode(errWriter{}, 0)
		require.ErrorContains(t, err, "write failed", mnem)
	}

	for _, f := range family3r {
		check(f.mnem, f.ctor(lreg(t, 12), lreg(t, 13), lreg(t, 14)))
	}

	for _, f := range family2r {
		check(f.mnem, f.ctor(lreg(t, 12), lreg(t, 13)))
	}
}
