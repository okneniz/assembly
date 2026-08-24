package alias

import (
	"encoding/binary"
	"testing"

	"github.com/okneniz/parsec/bytes"
	"github.com/stretchr/testify/require"

	arch "github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/asm"
	"github.com/okneniz/assembly/disasm"
)

func assembleOne(t *testing.T, src string, addr uint64) uint32 {
	t.Helper()
	res, errs := asm.Assemble(src, addr, NewASMBackend())
	require.Empty(t, errs, "assemble %q", src)
	require.NotEmpty(t, res.Sections, "assemble %q", src)
	require.Len(t, res.Sections[0].Data, 4, "assemble %q: bad output", src)
	return binary.LittleEndian.Uint32(res.Sections[0].Data)
}

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
	}
	for _, c := range cases {
		got := assembleOne(t, c.src, 0)
		require.Equal(t, c.word, got, "case %q", c.src)
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
		word := assembleOne(t, src, 0)
		insts, err := arch.Parse(0)(bytes.Buffer(binary.LittleEndian.AppendUint32(nil, word)))
		require.NoError(t, err)
		require.Len(t, insts, 1, "%q → %#08x: nothing decoded", src, word)
		text := insts[0].ObjDump(disasm.DefaultViewCtx())
		res, errs := asm.Assemble(text, insts[0].Addr(), NewASMBackend())
		require.Empty(t, errs, "%q → %#08x → %q: re-assemble", src, word, text)
		got := binary.LittleEndian.Uint32(res.Sections[0].Data)
		require.Equal(t, word, got, "%q → %#08x → %q", src, word, text)
	}
}
