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
		{"mul.w", NewMulW},
		{"mulh.w", NewMulhW},
		{"mulh.wu", NewMulhWu},
		{"mul.d", NewMulD},
		{"mulh.d", NewMulhD},
		{"mulh.du", NewMulhDu},
		{"mulw.d.w", NewMulwDW},
		{"mulw.d.wu", NewMulwDWu},
		{"div.w", NewDivW},
		{"mod.w", NewModW},
		{"div.wu", NewDivWu},
		{"mod.wu", NewModWu},
		{"div.d", NewDivD},
		{"mod.d", NewModD},
		{"div.du", NewDivDu},
		{"mod.du", NewModDu},
	}

	family2r := []struct {
		mnem string
		ctor func(rd, rj Reg) Instr
	}{
		{"ext.w.b", NewExtWB},
		{"ext.w.h", NewExtWH},
		{"rdtimel.w", NewRdtimelW},
		{"rdtimeh.w", NewRdtimehW},
		{"rdtime.d", NewRdtimeD},
		{"cpucfg", NewCpucfg},
		{"clo.w", NewCloW},
		{"clz.w", NewClzW},
		{"cto.w", NewCtoW},
		{"ctz.w", NewCtzW},
		{"clo.d", NewCloD},
		{"clz.d", NewClzD},
		{"cto.d", NewCtoD},
		{"ctz.d", NewCtzD},
		{"revb.2h", NewRevb2H},
		{"revb.4h", NewRevb4H},
		{"revb.2w", NewRevb2W},
		{"revb.d", NewRevbD},
		{"revh.2w", NewRevh2W},
		{"revh.d", NewRevhD},
		{"bitrev.4b", NewRevbit4B},
		{"bitrev.w", NewRevbitW},
		{"bitrev.8b", NewRevbit8B},
		{"bitrev.d", NewRevbitD},
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
