package alias

// The fragment oracle: a fragment resolved at an address is byte-identical
// to the full assembler's output for the same text at the same address -
// the two producers speak the same grammar and the same vocabulary. Plus
// the host-symbol path (a fragment calling a program label through the
// unit output) and the real rejections.

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	arch "github.com/okneniz/assembly/arch/arm64"
	asm "github.com/okneniz/assembly/asm"
	"github.com/okneniz/assembly/unit"
)

// fragBytes resolves the fragment at base and encodes it.
func fragBytes(t *testing.T, src string, base uint64, resolve func(string) (uint64, bool)) []byte {
	t.Helper()

	f, errs := Fragment(src)
	require.Empty(t, errs, src)

	rs, err := f.Resolve(unit.NewCtx(base, resolve))
	require.NoError(t, err, src)

	var buf bytes.Buffer
	for _, r := range rs {
		_, err := r.Encode(&buf)
		require.NoError(t, err, src)
	}

	return buf.Bytes()
}

func TestFragmentOracle(t *testing.T) {
	// data-only: the fragment text (aliases and numeric locals included);
	// the oracle is the full alias.Assemble of the same text
	tests := []string{
		"add w0, w1, w2",
		"cmp w0, #1",  // the alias family
		"mov x0, x1",  // the mov alias
		"cset x0, eq", // the cset alias
		"nop",
		"nop\nnop\nnop",
		"hlt #0xf000",    // the semihosting exit of an inline-asm body
		"hvc #0\nsmc #0", // the monitor calls of PSCI
		"1: nop\nb 1b",   // a numeric backward branch
		"b 1f\nnop\n1:",  // a numeric forward branch
		"adr x0, .",      // the fragment's own address
		"1: ldaxr w8, [x0]\nadd w8, w8, w1\nstlxr w9, w8, [x0]\ncbnz w9, 1b", // the atomic idiom
	}

	for _, src := range tests {
		res, errs := Assemble(src, 0x10000)
		require.Empty(t, errs, src)
		code := textOf(t, res)

		require.Equal(t, code, fragBytes(t, src, 0x10000, nil), src)
		require.Equal(t, len(code), func() int {
			f, ferrs := Fragment(src)
			require.Empty(t, ferrs, src)
			return f.Size()
		}())
	}
}

func TestFragmentHostSymbol(t *testing.T) {
	// asm("bl foo") inside a unit program: foo is a label of the HOST,
	// resolved against the unit's own layout (a nop after the fragment)
	u := unit.New()
	u.Label("start")
	u.Sym(unit.Pos{}, mustFragment(t, "bl foo"))
	u.Label("foo")
	u.Instr(unit.Pos{}, mustNop(t), nil)

	f := u.Resolve(func(text, data, dataMem int) (uint64, uint64) {
		require.Equal(t, 8, text)
		return 0x1000, 0x1000
	})
	require.Empty(t, f.Errs)

	code, codeErr := f.EncodeText()
	require.NoError(t, codeErr)

	// the oracle: the same text as a full source with foo right after
	res, errs := Assemble("bl foo\nfoo:\nnop", 0x1000)
	require.Empty(t, errs)
	require.Equal(t, textOf(t, res)[:4], code[:4])
}

func TestFragmentBadLocalAtResolve(t *testing.T) {
	// "1f" without a definition parses (the placeholder resolves it) and
	// fails only at resolve, carrying the fragment line
	f, errs := Fragment("nop\nb 1f")
	require.Empty(t, errs)

	_, err := f.Resolve(unit.NewCtx(0x1000, nil))
	require.ErrorContains(t, err, "line 2:")
}

func TestFragmentRejectsReal(t *testing.T) {
	// data-only: the source and the message fragment of the rejection
	tests := []struct {
		src string
		msg string
	}{
		{"ldr x0, =5", "literal pools"}, // the real PoolUser path
		{".word 5", "directive"},
		{"name: nop", "named label"},
		{"xyzzy w0", "unknown mnemonic"},
	}

	for _, tt := range tests {
		f, errs := Fragment(tt.src)
		require.Nil(t, f, tt.src)
		require.NotEmpty(t, errs, tt.src)
		require.ErrorContains(t, errs[0], tt.msg, tt.src)
	}
}

// textOf is the .text section data of a full assembly (the oracle bytes).
func textOf(t *testing.T, res *asm.Result) []byte {
	t.Helper()

	for _, s := range res.Sections {
		if s.Name == ".text" {
			return s.Data
		}
	}

	t.Fatal("no .text section in the oracle assembly")
	return nil
}

// mustFragment parses or fails the test.
func mustFragment(t *testing.T, src string) *asm.Fragment {
	t.Helper()

	f, errs := Fragment(src)
	require.Empty(t, errs)
	return f
}

// mustNop builds the nop through the arch Builder (the producer's path).
func mustNop(t *testing.T) arch.Instr {
	t.Helper()
	return arch.Builder{}.Nop()
}
