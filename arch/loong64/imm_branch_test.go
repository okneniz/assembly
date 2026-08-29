package loong64

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

// imm12v - a validated si12; a validation error fails the test.
func imm12v(t *testing.T, v int64) Imm12 {
	t.Helper()

	i, err := NewImm12(v)
	require.NoError(t, err)

	return i
}

// imm20v - a validated si20; a validation error fails the test.
func imm20v(t *testing.T, v int64) Imm20 {
	t.Helper()

	i, err := NewImm20(v)
	require.NoError(t, err)

	return i
}

func TestAddiWCtor(t *testing.T) {
	// llvm-mc-verified: addi.w $t0, $t1, -16.
	require.Equal(
		t,
		uint32(0x02bfc1ac),
		ctorWord(t, NewAddiW(lreg(t, 12), lreg(t, 13), imm12v(t, -16))),
	)
}

func TestAddiWDecodeEncode(t *testing.T) {
	in := decodeOne(0x02bfc1ac, 0x90000000)

	x, ok := in.(AddiW)
	require.True(t, ok, "type = %T, want AddiW", in)
	require.Equal(t, "addi.w $t0, $t1, -16", x.ObjDump(disasm.DefaultViewCtx()))

	// The negative immediate round-trips through the sign-extended field.
	require.Equal(t, uint32(0x02bfc1ac), ctorWord(t, x))
}

func TestLu12iWCtorDecode(t *testing.T) {
	// llvm-mc-verified: lu12i.w $t0, -1.
	in := NewLu12iW(lreg(t, 12), imm20v(t, -1))
	require.Equal(t, uint32(0x15ffffec), ctorWord(t, in))

	x, ok := decodeOne(0x15ffffec, 0).(Lu12iW)
	require.True(t, ok, "type = %T, want Lu12iW", x)
	require.Equal(t, "lu12i.w $t0, -1", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, int64(-1), x.imm.val)
}

func TestBeqCtorDecode(t *testing.T) {
	// llvm-mc-verified: beq $t1, $t0, 8 (at pc 0; the manual order prints
	// rj first). The constructor takes (rj, rd, target).
	in := NewBeq(lreg(t, 13), lreg(t, 12), 8)
	require.Equal(t, uint32(0x580009ac), ctorWord(t, in))

	x, ok := decodeOne(0x580009ac, 0).(Beq)
	require.True(t, ok, "type = %T, want Beq", x)
	require.Equal(t, "beq $t1, $t0, 8", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, int64(8), x.off.val)

	// The off field is the byte offset itself: the same word decodes
	// identically at any pc, and Encode is pc-independent.
	y, ok2 := decodeOne(0x580009ac, 0x90000000).(Beq)
	require.True(t, ok2, "type = %T, want Beq", y)
	require.Equal(t, int64(8), y.off.val)
	require.Equal(t, uint32(0x580009ac), ctorWord(t, y))
}

func TestBeqEncodeErrors(t *testing.T) {
	in := NewBeq(lreg(t, 13), lreg(t, 12), 6)

	_, err := in.Encode(errWriter{}, 0)
	require.ErrorContains(t, err, "not word-aligned")

	// The signed 16-bit word range is +-128 KiB.
	in = NewBeq(lreg(t, 13), lreg(t, 12), 1<<18)
	_, err = in.Encode(errWriter{}, 0)
	require.ErrorContains(t, err, "does not fit")
}

func TestBeqzCtorDecode(t *testing.T) {
	// llvm-mc-verified: beqz $t1, 8.
	in := NewBeqz(lreg(t, 13), 8)
	require.Equal(t, uint32(0x400009a0), ctorWord(t, in))

	x, ok := decodeOne(0x400009a0, 0).(Beqz)
	require.True(t, ok, "type = %T, want Beqz", x)
	require.Equal(t, "beqz $t1, 8", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, int64(8), x.off.val)

	// The split offs21 field: a negative offset round-trips; the off
	// field is the offset itself, identical at any pc.
	y, ok := decodeOne(0x400045a0, 0x1000).(Beqz)
	require.True(t, ok, "type = %T, want Beqz", y)
	require.Equal(t, int64(0x44), y.off.val)
	require.Equal(t, uint32(0x400045a0), ctorWord(t, y))
}

func TestBCtorDecode(t *testing.T) {
	// llvm-mc-verified: b 8 and b -8 (at pc 0).
	require.Equal(t, uint32(0x50000800), ctorWord(t, NewB(8)))
	require.Equal(t, uint32(0x53fffbff), ctorWord(t, NewB(-8)))

	x, ok := decodeOne(0x53fffbff, 0).(B)
	require.True(t, ok, "type = %T, want B", x)
	require.Equal(t, "b -8", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, int64(-8), x.off.val)
}

func TestJirlCtorDecode(t *testing.T) {
	// llvm-mc-verified: jirl $t0, $t1, 4.
	in := NewJirl(lreg(t, 12), lreg(t, 13), 4)
	require.Equal(t, uint32(0x4c0005ac), ctorWord(t, in))

	x, ok := decodeOne(0x4c0005ac, 0).(Jirl)
	require.True(t, ok, "type = %T, want Jirl", x)
	require.Equal(t, "jirl $t0, $t1, 4", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, int64(4), x.off.val)

	_, err := NewJirl(lreg(t, 0), lreg(t, 1), 3).Encode(errWriter{}, 0)
	require.ErrorContains(t, err, "not word-aligned")
}

func TestTemplatesJSONEncodeError(t *testing.T) {
	for _, tc := range []struct {
		mnem string
		in   Instr
	}{
		{"addi.w", NewAddiW(lreg(t, 12), lreg(t, 13), imm12v(t, -16))},
		{"lu12i.w", NewLu12iW(lreg(t, 12), imm20v(t, -1))},
		{"beq", NewBeq(lreg(t, 13), lreg(t, 12), 8)},
		{"beqz", NewBeqz(lreg(t, 13), 8)},
		{"b", NewB(8)},
		{"jirl", NewJirl(lreg(t, 12), lreg(t, 13), 4)},
	} {
		b, err := tc.in.MarshalJSON()
		require.NoError(t, err, tc.mnem)

		var dto map[string]any
		require.NoError(t, json.Unmarshal(b, &dto), tc.mnem)
		require.Equal(t, tc.mnem, dto["mnemonic"], tc.mnem)

		_, err = tc.in.Encode(errWriter{}, 0)
		require.ErrorContains(t, err, "write failed", tc.mnem)
	}
}
