package objdump

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseLine(t *testing.T) {
	cases := []struct {
		name    string
		line    string
		addr    uint64
		isInstr bool
	}{
		// Mach-O (ARM64): bytes separated by spaces.
		{
			"macho bytes4",
			"100003f64: fd 7b bf a9 stp x29, x30, [sp, #-0x10]!",
			0x100003f64,
			true,
		},
		{
			"macho bytes2 bare",
			"100003f68: 02 88",
			0x100003f68,
			true,
		},
		// ELF (RISC-V): hex word.
		{
			"elf word8",
			"a: ff010113 addi sp, -16",
			0xa,
			true,
		},
		{
			"elf word4",
			"1e: 1305 addi a0,a0,1",
			0x1e,
			true,
		},
		// A 16-bit compressed instruction as bytes (llvm does not print it
		// this way on Mach-O, but 2-byte fields are valid on their own).
		{
			"bytes2 rvc",
			"42: 02 88 c.eqz a0",
			0x42,
			true,
		},
		// Rejections: not instructions.
		{
			"header",
			"hello-world:",
			0,
			false,
		},
		{
			"section hdr",
			"Disassembly of section __TEXT,__text:",
			0,
			false,
		},
		{
			"symbol",
			"0000000100003f64 <_main>:",
			0,
			false,
		},
		{
			"3 bytes",
			"1000: fd 7b bf add x0, x0",
			0,
			false,
		},
		{
			"5 bytes",
			"1000: fd 7b bf a9 cc add x0, x0",
			0,
			false,
		},
		{
			"5-hex word",
			"1000: abcde addi a0, ",
			0,
			false,
		},
		{
			"10-hex word",
			"1000: ff010113aa addi a0, ",
			0,
			false,
		},
		{
			"empty",
			"",
			0,
			false,
		},
		{
			"no colon",
			"1000 fd 7b bf a9",
			0,
			false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			norm := Normalize(c.line)
			addr, ok := ParseLine(norm)
			require.Equal(t, c.isInstr, ok, "case %q: ParseLine(%q) recognized", c.name, norm)
			if ok {
				require.Equal(t, c.addr, addr, "case %q: ParseLine(%q) addr", c.name, norm)
			}
		})
	}
}

func TestStripComments(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{
			"adrp x27, 352 ; 0x100161000",
			"adrp x27, 352",
		},
		{
			"bl 0x100003f7c <_foo+0x14>",
			"bl 0x100003f7c",
		},
		{
			"add x0, x1, #4",
			"add x0, x1, #4",
		},
	}
	for _, c := range cases {
		require.Equal(t, c.want, StripComments(c.in), "StripComments(%q)", c.in)
	}
}

func TestNormalize(t *testing.T) {
	require.Equal(t, "a b c", Normalize("  a\t\tb   c\n"))
}

func TestParseByAddr(t *testing.T) {
	out := `hello-world:
(__TEXT,__text) section
100003f64: fd 7b bf a9 stp x29, x30, [sp, #-0x10]!
100003f68: fd 03 00 91 add x29, sp, #0
100003f6c: 13 04 00 94 bl 0x100003f7c <_foo+0x14>`

	m := ParseByAddr(out)
	require.Len(t, m, 3, "ParseByAddr: %v", m)
	s, ok := m[0x100003f6c]
	require.True(t, ok, "bad line at 0x100003f6c: %q", s)
	require.Equal(
		t,
		"100003f6c: 13 04 00 94 bl 0x100003f7c <_foo+0x14>",
		Normalize(s),
		"line at 0x100003f6c",
	)
}
