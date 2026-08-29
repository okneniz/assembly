package loong64

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

// code15v - a validated ui15 code; a validation error fails the test.
func code15v(t *testing.T, v int64) Code15 {
	t.Helper()

	c, err := NewCode15(v)
	require.NoError(t, err)

	return c
}

// TestBranchScaledOffsetErrors - the encPs2 error branch of every
// word-scaled offset instruction (alignment, range).
func TestBranchScaledOffsetErrors(t *testing.T) {
	branches := []struct {
		mnem string
		in   Instr
	}{
		{"bne", NewBne(lreg(t, 13), lreg(t, 12), 1<<18)},
		{"blt", NewBlt(lreg(t, 13), lreg(t, 12), 1<<18)},
		{"bge", NewBge(lreg(t, 13), lreg(t, 12), 1<<18)},
		{"bltu", NewBltu(lreg(t, 13), lreg(t, 12), 1<<18)},
		{"bgeu", NewBgeu(lreg(t, 13), lreg(t, 12), 1<<18)},
		{"bnez", NewBnez(lreg(t, 13), 1<<23)},
		{"bl", NewBl(1 << 28)},
	}
	// The branch offsets are raw byte offsets now: a span beyond the
	// field width does not fit (the ldptr/stptr/ll/sc rows are absent -
	// their offsets are role-validated up front and cannot fail here).
	for _, tc := range branches {
		_, err := tc.in.Encode(errWriter{}, 0)
		require.ErrorContains(t, err, "does not fit", tc.mnem)
	}

	// jirl's offset is rj-relative: the raw constructor takes any int64,
	// the range check is genuinely dynamic.
	_, err := NewJirl(lreg(t, 12), lreg(t, 13), 1<<18).Encode(errWriter{}, 0)
	require.ErrorContains(t, err, "does not fit")
}

// TestDecodeAliases - the llvm-printed alias forms (verified word for word
// by TestDiffAgainstLLVM; these pin the exact branch of each alias).
func TestDecodeAliases(t *testing.T) {
	for _, tc := range []struct {
		w    uint32
		want string
	}{
		// blt rj, $zero -> bltz rj; blt $zero, rd -> bgtz rd (rd-check first).
		{0x600001a0, "bltz $t1, 0"},
		{0x6000000c, "bgtz $t0, 0"},
		{0x60000000, "bltz $zero, 0"},
		// bge $zero, rd -> blez rd; bge rj, $zero -> bgez rj (rj-check first).
		{0x640001a0, "bgez $t1, 0"},
		{0x6400000c, "blez $t0, 0"},
		{0x64000000, "blez $zero, 0"},
		// jirl $zero, rj, 0 -> jr rj (and ret for $ra).
		{0x4c000000 | 13<<5, "jr $t1"},
		{0x4c000020, "ret"},
		// rdtimel.w $zero, rj -> rdcntid.w rj; rdtimel.w rd, $zero -> rdcntvl.w.
		{0x000061a0, "rdcntid.w $t1"},
		{0x0000600c, "rdcntvl.w $t0"},
		{0x0000640c, "rdcntvh.w $t0"},
		// or rd, rj, $zero -> move; andi $zero, $zero, 0 -> nop.
		{0x001501ac, "move $t0, $t1"},
		{0x03400000, "nop"},
	} {
		require.Equal(
			t,
			tc.want,
			decodeOne(tc.w, 0).ObjDump(disasm.DefaultViewCtx()),
			"%#x",
			tc.w,
		)
	}
}

// TestLateFilesJSONEncodeError - the JSON and write-error paths of the
// instructions added after the slice family tests were written.
func TestLateFilesJSONEncodeError(t *testing.T) {
	for _, tc := range []struct {
		mnem string
		in   Instr
	}{
		{"ld.bu", NewLdBu(lreg(t, 12), lreg(t, 13), imm12v(t, 8))},
		{"ld.hu", NewLdHu(lreg(t, 12), lreg(t, 13), imm12v(t, 8))},
		{"ldx.bu", NewLdxBu(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"ldx.hu", NewLdxHu(lreg(t, 12), lreg(t, 13), lreg(t, 14))},
		{"dbcl", NewDbcl(code15v(t, 1))},
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
