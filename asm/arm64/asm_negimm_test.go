package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	asm "github.com/okneniz/assembly/asm"
)

// TestAddSubNegImmWords - the LLVM negative-immediate canon of the
// add/sub family: a negative imm12 flips the base (add x0, x1, #-1 is
// the sub x0, x1, #1 alias, cmp x0, #-1 the cmn one), a value beyond
// 12 bits folds into the imm12, lsl #12 form. Every word is
// clang-pinned.
func TestAddSubNegImmWords(t *testing.T) {
	cases := []struct {
		src  string
		word uint32
	}{
		{"cmp x0, #-1", 0xb100041f},
		{"cmp w0, #-1", 0x3100041f},
		{"cmp w0, #-5", 0x3100141f},
		{"cmn x0, #-1", 0xf100041f},
		{"cmp x0, #-4096", 0xb140041f},
		{"cmp x0, #-8192", 0xb140081f},
		{"subs x9, x10, #-4095", 0xb13ffd49},
		{"add x0, x1, #-1", 0xd1000420},
		{"sub x0, x1, #-1", 0x91000420},
		{"adds x0, x1, #-1", 0xf1000420},
		{"adds w2, w3, #-7", 0x71001c62},
		{"subs x0, x1, #-1", 0xb1000420},
		{"add sp, sp, #-16", 0xd10043ff},
		{"cmp x0, #4096", 0xf140041f},
		{"add x0, x1, #4096", 0x91400420},
	}

	for _, c := range cases {
		require.Equal(t, c.word, armAssembleOne(t, c.src+"\n", 0), "case %q", c.src)
	}
}

// TestAddSubNegImmErrors - the unencodable negatives: beyond the folded
// 24-bit span or not a multiple of 4096 within it.
func TestAddSubNegImmErrors(t *testing.T) {
	cases := []string{
		"cmp x0, #-5000",
		"cmp x0, #-16777217",
		"add x0, x1, #-4097",
	}

	for _, src := range cases {
		_, errs := asm.Assemble(src+"\n", 0, New())
		require.NotEmpty(t, errs, "case %q must not assemble", src)
	}
}
