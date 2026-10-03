package alias

import (
	"encoding/binary"
	"testing"

	"github.com/okneniz/parsec"
	"github.com/okneniz/parsec/bytes"
	"github.com/stretchr/testify/require"

	arch "github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/asm"
	"github.com/okneniz/assembly/disasm"
)

// TestAliasWords checks the exact words (verified against the encodings
// of the base instructions: an alias has the same encoding as the base
// form).
func TestAliasWords(t *testing.T) {
	cases := []struct {
		src  string
		word uint32
	}{
		{
			"mov x0, #0x1",
			0xd2800020,
		}, // movz x0, #1
		{
			"cmp x0, #1",
			0xf100041f,
		}, // subs xzr, x0, #1
		{
			"neg x0, x1",
			0xcb0103e0,
		}, // sub x0, xzr, x1
		{
			"cset x0, eq",
			0x9a9f17e0,
		}, // csinc x0, xzr, xzr, ne
		{
			"mul x0, x1, x2",
			0x9b027c20,
		}, // madd x0, x1, x2, xzr
		{
			"mov x6, #-0x66660001",
			0x92acccc6,
		}, // movn x6, #0x6666, lsl #16
		{
			"mov x0, #-0x123400000001",
			0x92c24680,
		}, // movn x0, #0x1234, lsl #32
		{
			"mov w0, #-0x12340001",
			0x12a24680,
		}, // movn w0, #0x1234, lsl #16
	}
	for _, c := range cases {
		got := assembleOne(t, c.src)
		require.Equal(t, c.word, got, "case %q", c.src)
	}
}

// TestMovNegImmRefusal — the mov #imm negatives that no single
// instruction encodes (not MOVN-shaped, not a logical immediate) are
// refused, and the 64-bit-only lanes never reach a W form.
func TestMovNegImmRefusal(t *testing.T) {
	for _, src := range []string{
		"mov x6, #-0x66660000",  // not ~(imm16 << 16hw)
		"mov w0, #-0x100000001", // hw>=2 in a W form
	} {
		res, errs := asm.Assemble(src, 0, NewASMBackend())
		require.NotEmpty(t, errs, "case %q must not assemble", src)
		require.Empty(t, res.Sections, "case %q", src)
	}
}

// TestAliasRoundTrip runs each alias through assemble -> decode ->
// decoder text -> assemble -> same word (an outer-level self-verify:
// the constructors agree with the decoder formatters).
func TestAliasRoundTrip(t *testing.T) {
	cases := []string{
		"cmp x0, #1", "cmn x0, #1", "cmp x1, x2, lsl #4", "cmp x1, x2, uxtx #1",
		"neg x0, x1", "negs w0, w1", "neg x0, x1, asr #8",
		"tst x0, #0xf", "tst x0, x1", "tst x0, x1, lsl #4",
		"mvn x0, x1", "mvn w0, w1, ror #8",
		"mov x0, x1", "mov x0, #0x42", "mov x0, #0x12340000", "mov w0, #-1",
		"mul x0, x1, x2", "mneg x0, x1, x2",
		"cset x0, eq", "csetm x0, ne",
		"cinc x0, x1, mi", "cinv x0, x1, pl", "cneg x0, x1, vs",
		"sxtb x0, w1", "sxth x0, w1", "sxtw x0, w1",
		"ubfiz x0, x1, #4, #8", "ubfx x0, x1, #4, #8",
		"sbfiz x0, x1, #4, #8", "sbfx x0, x1, #4, #8",
	}
	for _, src := range cases {
		word := assembleOne(t, src)
		insts, err := arch.MakeDecoder()(
			parsec.Stateless{},
			bytes.Buffer(binary.LittleEndian.AppendUint32(nil, word)),
		)
		require.NoError(t, err)
		require.Len(t, insts, 1, "%q → %#08x: nothing decoded", src, word)
		text := insts[0].ObjDump(disasm.DefaultViewCtx())
		res, errs := asm.Assemble(text, 0, NewASMBackend())
		require.Empty(t, errs, "%q → %#08x → %q: re-assemble", src, word, text)
		got := binary.LittleEndian.Uint32(res.Sections[0].Data)
		require.Equal(t, word, got, "%q → %#08x → %q", src, word, text)
	}
}

// TestMovFromVectorAlias - the MOV (from vector) input alias of UMOV:
// llvm accepts both spellings for the filling sizes (s into w, d into x)
// and prints mov; mov wd, vn.b[n]/vn.h[n] is not a spelling there, and
// the umov width must match the element. Words pinned against llvm-mc
// (docker assembly-tests, 2026-09-08).
func TestMovFromVectorAlias(t *testing.T) {
	for _, c := range []struct {
		src  string
		word uint32
	}{
		{"mov x0, v1.d[1]", 0x4E183C20},
		{"mov w0, v1.s[1]", 0x0E0C3C20},
		{"umov x0, v1.d[1]", 0x4E183C20},
		{"umov w0, v1.s[1]", 0x0E0C3C20},
	} {
		require.Equal(t, c.word, assembleOne(t, c.src), "case %q", c.src)
	}

	for _, src := range []string{
		"mov w0, v1.b[3]",
		"mov w0, v1.h[7]",
		"mov x0, v1.s[1]",
		"umov w0, v1.d[1]",
		"umov x0, v1.s[1]",
	} {
		res, errs := asm.Assemble(src, 0, NewASMBackend())
		require.NotEmpty(t, errs, "case %q must not assemble", src)
		require.Empty(t, res.Sections, "case %q", src)
	}
}

// TestDupScalarAlias - the scalar DUP alias mov <b|h|s|d>n, vn.sz[idx]:
// the destination view class must match the source element letter, the
// index stays inside the lane count. Words pinned against clang
// (arm64-none-elf, 2026-10-02).
func TestDupScalarAlias(t *testing.T) {
	for _, c := range []struct {
		src  string
		word uint32
	}{
		{"mov s9, v19.s[0]", 0x5e040669},
		{"mov b9, v19.b[3]", 0x5e070669},
		{"mov h9, v19.h[7]", 0x5e1e0669},
		{"mov d9, v19.d[1]", 0x5e180669},
	} {
		require.Equal(t, c.word, assembleOne(t, c.src), "case %q", c.src)
	}

	for _, src := range []string{
		"mov s9, v19.h[0]", // the view class must match the element
		"mov s9, v19.s[4]", // the lane index out of the .s count
		"mov d9, v19.d[2]",
		"mov v9.s, v19.s[0]", // the vector dest is not a scalar view
	} {
		res, errs := asm.Assemble(src, 0, NewASMBackend())
		require.NotEmpty(t, errs, "case %q must not assemble", src)
		require.Empty(t, res.Sections, "case %q", src)
	}
}

// TestAliasBitfieldWords - the width-fixed alias forms of the bitfield
// round, the words llvm-mc emits (clang 19.1.7).
func TestAliasBitfieldWords(t *testing.T) {
	cases := []struct {
		src  string
		word uint32
	}{
		{"sxtb w0, w1", 0x13001c20},         // sbfm w0, w1, #0, #7
		{"sxth w0, w1", 0x13003c20},         // sbfm w0, w1, #0, #15
		{"sxtb x0, w1", 0x93401c20},         // sbfm x0, x1, #0, #7
		{"sxtw x0, w1", 0x93407c20},         // sbfm x0, x1, #0, #31
		{"uxtb w0, w1", 0x53001c20},         // ubfm w0, w1, #0, #7
		{"uxtb x0, w1", 0x53001c20},         // clang canonicalizes the x spelling
		{"uxth w0, w1", 0x53003c20},         // ubfm w0, w1, #0, #15
		{"ror w0, w1, #7", 0x13811c20},      // extr w0, w1, w1, #7
		{"ror x0, x1, #7", 0x93c11c20},      // extr x0, x1, x1, #7
		{"extr w0, w1, w2, #3", 0x13820c20}, // extr w0, w1, w2, #3
		{"extr wzr, w2, w2, #31", 0x13827c5f},
		{"lsl w0, w1, #5", 0x531b6820},
		{"lsr w0, w1, #5", 0x53057c20},
		{"asr w0, w1, #31", 0x131f7c20},
	}

	for _, c := range cases {
		require.Equal(t, c.word, assembleOne(t, c.src), "case %q", c.src)
	}

	for _, src := range []string{
		"sxtw w0, w1",     // sxtw has no 32-bit form
		"uxtw x0, w1",     // no such standalone alias
		"ror w0, w1, #32", // the W-form lsb range is 0..31
		"extr w0, w1, w2, #32",
		"lsl w0, w1, #32", // the W-form shift range is 0..31
		"lsr w0, w1, #32",
		"asr w0, w1, #32",
	} {
		res, errs := asm.Assemble(src, 0, NewASMBackend())
		require.NotEmpty(t, errs, "case %q must not assemble", src)
		require.Empty(t, res.Sections, "case %q", src)
	}
}

func assembleOne(t *testing.T, src string) uint32 {
	t.Helper()
	res, errs := asm.Assemble(src, 0, NewASMBackend())
	require.Empty(t, errs, "assemble %q", src)
	require.NotEmpty(t, res.Sections, "assemble %q", src)
	require.Len(t, res.Sections[0].Data, 4, "assemble %q: bad output", src)
	return binary.LittleEndian.Uint32(res.Sections[0].Data)
}
