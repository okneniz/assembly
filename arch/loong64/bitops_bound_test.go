package loong64

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

// uimm2v - a validated ui2; a validation error fails the test.
func uimm2v(t *testing.T, v int64) UImm2 {
	t.Helper()

	i, err := NewUImm2(v)
	require.NoError(t, err)

	return i
}

// uimm3v - a validated ui3; a validation error fails the test.
func uimm3v(t *testing.T, v int64) UImm3 {
	t.Helper()

	i, err := NewUImm3(v)
	require.NoError(t, err)

	return i
}

// uimm5v - a validated ui5; a validation error fails the test.
func uimm5v(t *testing.T, v int64) UImm5 {
	t.Helper()

	i, err := NewUImm5(v)
	require.NoError(t, err)

	return i
}

// uimm6v - a validated ui6; a validation error fails the test.
func uimm6v(t *testing.T, v int64) UImm6 {
	t.Helper()

	i, err := NewUImm6(v)
	require.NoError(t, err)

	return i
}

// shift3v - a validated alsl shift amount (1..4); an error fails the test.
func shift3v(t *testing.T, v int64) Shift3 {
	t.Helper()

	s, err := NewShift3(v)
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
		{"bstrins.w", NewBstrinsW(lreg(t, 12), lreg(t, 13), uimm5v(t, 5), uimm5v(t, 3))},
		{"bstrpick.w", NewBstrpickW(lreg(t, 12), lreg(t, 13), uimm5v(t, 5), uimm5v(t, 3))},
		{"bstrins.d", NewBstrinsD(lreg(t, 12), lreg(t, 13), uimm6v(t, 5), uimm6v(t, 3))},
		{"bstrpick.d", NewBstrpickD(lreg(t, 12), lreg(t, 13), uimm6v(t, 5), uimm6v(t, 3))},
		{"alsl.w", NewAlslW(lreg(t, 12), lreg(t, 13), lreg(t, 14), shift3v(t, 3))},
		{"alsl.wu", NewAlslWu(lreg(t, 12), lreg(t, 13), lreg(t, 14), shift3v(t, 3))},
		{"alsl.d", NewAlslD(lreg(t, 12), lreg(t, 13), lreg(t, 14), shift3v(t, 3))},
		{"bytepick.w", NewBytepickW(lreg(t, 12), lreg(t, 13), lreg(t, 14), uimm2v(t, 3))},
		{"bytepick.d", NewBytepickD(lreg(t, 12), lreg(t, 13), lreg(t, 14), uimm3v(t, 3))},
		{"crc.w.b.w", NewCrcWBW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"crc.w.h.w", NewCrcWHW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"crc.w.w.w", NewCrcWWW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"crc.w.d.w", NewCrcWDW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"crcc.w.b.w", NewCrccWBW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"crcc.w.h.w", NewCrccWHW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"crcc.w.w.w", NewCrccWWW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"crcc.w.d.w", NewCrccWDW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"asrtle.d", NewAsrtleD(lreg(t, 13), lreg(t, 14))},
		{"asrtgt.d", NewAsrtgtD(lreg(t, 13), lreg(t, 14))},
		{"ldgt.b", NewLdgtB(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"ldgt.h", NewLdgtH(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"ldgt.w", NewLdgtW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"ldgt.d", NewLdgtD(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"ldle.b", NewLdleB(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"ldle.h", NewLdleH(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"ldle.w", NewLdleW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"ldle.d", NewLdleD(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"stgt.b", NewStgtB(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"stgt.h", NewStgtH(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"stgt.w", NewStgtW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"stgt.d", NewStgtD(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"stle.b", NewStleB(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"stle.h", NewStleH(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"stle.w", NewStleW(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"stle.d", NewStleD(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
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
