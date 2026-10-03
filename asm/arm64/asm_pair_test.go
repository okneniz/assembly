package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	arch "github.com/okneniz/assembly/arch/arm64"
	asm "github.com/okneniz/assembly/asm"
	"github.com/okneniz/assembly/disasm"
)

// TestPairAssembleWords - the ldp/stp register classes (x/w/s/d/q) in the
// three addressing forms plus ldpsw; every word is clang-pinned
// (cross-checked on the host: clang -target aarch64-linux-gnu + objdump).
func TestPairAssembleWords(t *testing.T) {
	cases := []struct {
		src  string
		word uint32
	}{
		{"stp w0, w1, [x6, #16]", 0x290204c0},
		{"ldp w0, w1, [x6, #16]", 0x294204c0},
		{"stp w2, w3, [sp, #-32]!", 0x29bc0fe2},
		{"ldp w2, w3, [sp], #32", 0x28c40fe2},
		{"stp x0, x1, [x6, #16]", 0xa90104c0},
		{"ldp x0, x1, [x6, #16]", 0xa94104c0},
		{"stp x29, x30, [sp, #-16]!", 0xa9bf7bfd},
		{"ldp x29, x30, [sp], #16", 0xa8c17bfd},
		{"stp s0, s1, [x7, #-32]", 0x2d3c04e0},
		{"ldp s2, s3, [x6, #16]", 0x2d420cc2},
		{"stp s0, s1, [sp, #-32]!", 0x2dbc07e0},
		{"ldp s0, s1, [sp], #32", 0x2cc407e0},
		{"stp d0, d1, [x6, #16]", 0x6d0104c0},
		{"ldp d2, d3, [x7]", 0x6d400ce2},
		{"stp d8, d9, [sp, #-16]!", 0x6dbf27e8},
		{"ldp d8, d9, [sp], #16", 0x6cc127e8},
		{"stp q0, q1, [x6, #16]", 0xad0084c0},
		{"ldp q30, q31, [x6, #480]", 0xad4f7cde},
		{"stp q0, q1, [sp, #-32]!", 0xadbf07e0},
		{"ldp q0, q1, [sp], #32", 0xacc107e0},
		{"ldpsw x0, x1, [x6, #8]", 0x694104c0},
		{"ldpsw x2, x3, [sp, #-16]!", 0x69fe0fe2},
		{"ldpsw x4, x5, [sp], #16", 0x68c217e4},
		{"stp xzr, xzr, [sp, #-16]!", 0xa9bf7fff},
	}

	for _, c := range cases {
		require.Equal(t, c.word, armAssembleOne(t, c.src+"\n", 0), "case %q", c.src)
	}
}

// TestPairAssembleErrors - the refused shapes: mixed-width and non-pair
// register classes, offsets off the imm7 grid or out of its range, and
// fp registers where only x/w live (ldpsw).
func TestPairAssembleErrors(t *testing.T) {
	cases := []string{
		"stp q0, d1, [x6]",
		"stp x0, w1, [x6]",
		"stp v0, v1, [x6]",
		"stp b0, b1, [x6]",
		"stp q0, q1, [x6, #12]",
		"stp q0, q1, [x6, #1024]",
		"stp d0, d1, [x6, #1016]",
		"ldpsw q0, q1, [x6, #8]",
	}

	for _, src := range cases {
		_, errs := asm.Assemble(src+"\n", 0, New())
		require.NotEmpty(t, errs, "case %q must not assemble", src)
	}
}

// TestPairDecodeTexts - the decoded text of the s/d/q pair words and the
// ldpsw offset form (the pre/post ldpsw words encode right but decode as
// the offset form - the Ldpsw struct knows no writeback).
func TestPairDecodeTexts(t *testing.T) {
	cases := []struct {
		word uint32
		text string
	}{
		{0xad0084c0, "stp q0, q1, [x6, #0x10]"},
		{0xad4f7cde, "ldp q30, q31, [x6, #0x1e0]"},
		{0xadbf07e0, "stp q0, q1, [sp, #-0x20]!"},
		{0xacc107e0, "ldp q0, q1, [sp], #0x20"},
		{0x2d3c04e0, "stp s0, s1, [x7, #-0x20]"},
		{0x2d420cc2, "ldp s2, s3, [x6, #0x10]"},
		{0x6d0104c0, "stp d0, d1, [x6, #0x10]"},
		{0x6dbf27e8, "stp d8, d9, [sp, #-0x10]!"},
		{0x6cc127e8, "ldp d8, d9, [sp], #0x10"},
		{0x290204c0, "stp w0, w1, [x6, #0x10]"},
		{0x694104c0, "ldpsw x0, x1, [x6, #0x8]"},
	}

	for _, c := range cases {
		ins, err := arch.DecodeWord(c.word)
		require.NoError(t, err, "case %#x", c.word)
		require.Equal(t, c.text, ins.ObjDump(disasm.ViewCtxAt(0)), "case %#x", c.word)
	}
}
