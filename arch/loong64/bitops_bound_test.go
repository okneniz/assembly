package loong64

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

// uimm2v - a validated ui2; a validation error fails the test.
func uimm2v(t *testing.T, v int64) UImm2 {
	t.Helper()

	i, err := New().UImm2(v)
	require.NoError(t, err)

	return i
}

// uimm3v - a validated ui3; a validation error fails the test.
func uimm3v(t *testing.T, v int64) UImm3 {
	t.Helper()

	i, err := New().UImm3(v)
	require.NoError(t, err)

	return i
}

// uimm5v - a validated ui5; a validation error fails the test.
func uimm5v(t *testing.T, v int64) UImm5 {
	t.Helper()

	i, err := New().UImm5(v)
	require.NoError(t, err)

	return i
}

// uimm6v - a validated ui6; a validation error fails the test.
func uimm6v(t *testing.T, v int64) UImm6 {
	t.Helper()

	i, err := New().UImm6(v)
	require.NoError(t, err)

	return i
}

// shift3v - a validated alsl shift amount (1..4); an error fails the test.
func shift3v(t *testing.T, v int64) Shift3 {
	t.Helper()

	s, err := New().Shift3(v)
	require.NoError(t, err)

	return s
}

// TestBitopsBoundJSONEncodeError - the marshal and write-error paths of
// every bit-field, alsl, bytepick, crc, assert and bounds-check
// instruction (the rest is covered by the per-file tests).
func TestBitopsBoundJSONEncodeError(t *testing.T) {
	family := []struct {
		mnem string
		in   Instr
	}{
		{"bstrins.w", New().BstrinsW(lreg(t, 12), lreg(t, 13), uimm5v(t, 5), uimm5v(t, 3))},
		{"bstrpick.w", New().BstrpickW(lreg(t, 12), lreg(t, 13), uimm5v(t, 5), uimm5v(t, 3))},
		{"bstrins.d", New().BstrinsD(lreg(t, 12), lreg(t, 13), uimm6v(t, 5), uimm6v(t, 3))},
		{"bstrpick.d", New().BstrpickD(lreg(t, 12), lreg(t, 13), uimm6v(t, 5), uimm6v(t, 3))},
		{"alsl.w", New().AlslW(lreg(t, 12), lreg(t, 13), lreg(t, 14), shift3v(t, 3))},
		{"alsl.wu", New().AlslWu(lreg(t, 12), lreg(t, 13), lreg(t, 14), shift3v(t, 3))},
		{"alsl.d", New().AlslD(lreg(t, 12), lreg(t, 13), lreg(t, 14), shift3v(t, 3))},
		{"bytepick.w", New().BytepickW(lreg(t, 12), lreg(t, 13), lreg(t, 14), uimm2v(t, 3))},
		{"bytepick.d", New().BytepickD(lreg(t, 12), lreg(t, 13), lreg(t, 14), uimm3v(t, 3))},
		{"crc.w.b.w", New().CrcWBW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"crc.w.h.w", New().CrcWHW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"crc.w.w.w", New().CrcWWW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"crc.w.d.w", New().CrcWDW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"crcc.w.b.w", New().CrccWBW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"crcc.w.h.w", New().CrccWHW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"crcc.w.w.w", New().CrccWWW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"crcc.w.d.w", New().CrccWDW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"asrtle.d", New().AsrtleD(lreg(t, 13), lreg(t, 14))},
		{"asrtgt.d", New().AsrtgtD(lreg(t, 13), lreg(t, 14))},
		{"ldgt.b", New().LdgtB(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"ldgt.h", New().LdgtH(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"ldgt.w", New().LdgtW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"ldgt.d", New().LdgtD(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"ldle.b", New().LdleB(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"ldle.h", New().LdleH(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"ldle.w", New().LdleW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"ldle.d", New().LdleD(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"stgt.b", New().StgtB(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"stgt.h", New().StgtH(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"stgt.w", New().StgtW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"stgt.d", New().StgtD(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"stle.b", New().StleB(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"stle.h", New().StleH(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"stle.w", New().StleW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"stle.d", New().StleD(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
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
