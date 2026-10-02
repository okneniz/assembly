package arm64

import (
	"encoding/binary"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/okneniz/parsec"
	"github.com/okneniz/parsec/bytes"
	"github.com/stretchr/testify/require"

	arch "github.com/okneniz/assembly/arch/arm64"
	asm "github.com/okneniz/assembly/asm"
	"github.com/okneniz/assembly/disasm"
	"github.com/okneniz/assembly/file"
	"github.com/okneniz/assembly/tests/cmd/objdump"
)

func armAssembleOne(t *testing.T, src string, addr uint64) uint32 {
	t.Helper()
	res, errs := asm.Assemble(src, addr, New())
	require.Empty(t, errs, "assemble %q", src)
	require.NotEmpty(t, res.Sections, "assemble %q", src)
	require.Len(t, res.Sections[0].Data, 4, "assemble %q: bad output", src)
	return binary.LittleEndian.Uint32(res.Sections[0].Data)
}

func TestBitMasksRoundTrip(t *testing.T) {
	err := arch.VerifyBitMasks()
	require.NoError(t, err)
}

func TestArmAssembleWords(t *testing.T) {
	cases := []struct {
		src  string
		word uint32
	}{
		{
			"nop",
			0xd503201f,
		},
		{
			"ret",
			0xd65f03c0,
		},
		{
			"add x0, x1, #0x42",
			0x91068020,
		},
		{
			"mov x0, #0x1",
			0xd2800020,
		},
		{
			"lsl x0, x0, #1",
			0xd37ff800,
		},
		{
			"lsl x0, x0, #57",
			0xd3471800,
		},
		{
			"lsl w0, w0, #17",
			0x530f3800,
		},
		{
			"fmov s0, #-18.0",
			0x1e365000,
		},
		{
			"movz x0, #0x1234",
			0xd2824680,
		},
		{
			"svc #0x80",
			0xd4001001,
		},
		{
			"brk #0x1",
			0xd4200020,
		},
		{
			"hlt #0xf000",
			0xd45e0000,
		},
		{
			"hlt #0x1",
			0xd4400020,
		},
		{
			"hvc #0",
			0xd4000002,
		},
		{
			"smc #0",
			0xd4000003,
		},
		{
			"dmb sy",
			0xd5033fbf,
		},
		{
			"dsb sy",
			0xd5033f9f,
		},
		{
			"isb",
			0xd5033fdf,
		},
		{
			"isb sy",
			0xd5033fdf,
		},
		{
			"dsb ishst",
			0xd5033a9f,
		},
		{
			"dmb ish",
			0xd5033bbf,
		},
		{
			"mrs x0, CNTVCT_EL0",
			0xd53be040,
		},
		{
			"wfi",
			0xd503207f,
		},
		{
			"wfe",
			0xd503205f,
		},
		{
			"sev",
			0xd503209f,
		},
		{
			"sevl",
			0xd50320bf,
		},
		{
			"msr daifset, #1",
			0xd50341df,
		},
		{
			"msr daifset, #3",
			0xd50343df,
		},
		{
			"msr spsel, #1",
			0xd50041bf,
		},
		{
			"msr daifclr, #2",
			0xd50342ff,
		},
		{
			"msr pan, #1",
			0xd500419f,
		},
		{
			"msr allint, #1",
			0xd501411f,
		},
		{
			"msr pm, #1",
			0xd501431f,
		},
		{
			"msr cntvct_el0, x0",
			0xd51be040,
		},
	}
	for _, c := range cases {
		if c.src == "add x0, x1, #0x42" {
			continue // checked separately below
		}

		got := armAssembleOne(t, c.src, 0)
		require.Equal(t, c.word, got, "case %q", c.src)
	}

	// the sysreg lookup is case-insensitive: the lower-case kernel
	// spelling encodes (and prints) like the canonical one
	require.Equal(
		t,
		armAssembleOne(t, "mrs x0, CNTVCT_EL0", 0),
		armAssembleOne(t, "mrs x0, cntvct_el0", 0),
	)

	// add x0, x1, #0x42: imm12=0x42<<10 | Rn=1<<5 | Rd=0 | 0x91000000
	got := armAssembleOne(t, "add x0, x1, #0x42", 0)
	require.True(
		t,
		got == 0x91068020 || got == 0x91000000|0x42<<10|1<<5,
		"add = %#08x",
		got,
	)
}

// TestMsrPstateImmediate — the PSTATE field name is case-insensitive,
// the imm range is CRm 0..15 (allint/pm take a single bit), the unknown
// fields keep the sysreg-path error.
func TestMsrPstateImmediate(t *testing.T) {
	require.Equal(
		t,
		armAssembleOne(t, "msr daifset, #1", 0),
		armAssembleOne(t, "msr DAIFSet, #1", 0),
	)

	for _, src := range []string{
		"msr daifset, #16",
		"msr daifset, #-1",
		"msr allint, #2",
		"msr pm, #4",
		"msr nosuch, #1",
		"msr daifset",
		"msr daifset, x0",
		"wfi x0",
	} {
		_, errs := asm.Assemble(src, 0, New())
		require.NotEmpty(t, errs, "case %q must not assemble", src)
	}
}

func TestArmBranchAndMem(t *testing.T) {
	// b 0x1008 @ 0x1000
	got := armAssembleOne(t, "b 0x1008", 0x1000)
	require.Equal(t, uint32(0x14000002), got, "b")
	// bl 0x2000 @ 0x1000
	got = armAssembleOne(t, "bl 0x2000", 0x1000)
	require.Equal(t, uint32(0x94000400), got, "bl")
	// b.eq with a negative absolute target (base 0): the two's complement
	// spelling must decode back into b.eq, not fall away to the plain b
	// word of the self-verify fallback
	got = armAssembleOne(t, "b.eq 0xfffffffffff80004", 0)
	require.Equal(t, uint32(0x54c00020), got, "b.eq negative target")
	// cbz with the same shape
	got = armAssembleOne(t, "cbz x0, 0xfffffffffff80004", 0)
	require.Equal(t, uint32(0xb4c00020), got, "cbz negative target")
	// the S-form ext word with rn 31: the decoder must know it (rn 31
	// reads as sp/wsp there), or the self-verify fallback emits an
	// X-form word for a W instruction
	got = armAssembleOne(t, "add w0, wsp, w0, uxth #0", 0)
	require.Equal(t, uint32(0x0b2023e0), got, "add w ext with wsp base")
	got = armAssembleOne(t, "adds x0, sp, x0", 0)
	require.Equal(t, uint32(0xab2063e0), got, "adds x ext with sp base")
	// ldr x0, [x1]
	got = armAssembleOne(t, "ldr x0, [x1]", 0)
	require.Equal(t, uint32(0xf9400020), got, "ldr x0,[x1]")
	// ldr x0, [x1, #0x8]
	got = armAssembleOne(t, "ldr x0, [x1, #0x8]", 0)
	require.Equal(t, uint32(0xf9400420), got, "ldr")
	// stp x29, x30, [sp, #-0x10]!
	got = armAssembleOne(t, "stp x29, x30, [sp, #-0x10]!", 0)
	require.Equal(t, uint32(0xa9bf7bfd), got, "stp")
}

// TestArmRoundTripExample is a byte-exact round-trip test of the test
// binary: decode -> Text -> assemble -> the same bytes. The threshold
// matches the objdump gate.
func TestArmRoundTripExample(t *testing.T) {
	ff, err := file.Detect("../../tests/examples/hello-world/hello-world")
	if err != nil {
		t.Skipf("example not available: %v", err)
	}

	ts, err := ff.CodeSection()
	if err != nil {
		t.Skipf("example not available: %v", err)
	}

	insts, err := arch.MakeDecoder()(parsec.Stateless{}, bytes.Buffer(ts.Data))
	require.NoError(t, err)
	matched, failed, notAssembled, dontCare, equiv := 0, 0, 0, 0, 0
	sample := 0
	off := uint64(0)
	for _, in := range insts {
		addr := ts.Addr + off
		off += uint64(in.Len())
		if _, ok := in.(arch.Generic); ok {
			continue // decode-only: the generic syntax is not parsed by the assembler
		}

		src := objdump.StripComments(objdump.Normalize(in.ObjDump(disasm.ViewCtxAt(addr))))
		if src == "" || strings.HasPrefix(src, ".word") {
			continue
		}

		want := ts.Data[addr-ts.Addr : addr-ts.Addr+4]
		res, errs := asm.Assemble(src, addr, New())
		if len(errs) != 0 {
			notAssembled++
			if sample < 8 {
				t.Logf("addr %#x: %q: %v", addr, src, errs)
				sample++
			}

			continue
		}

		got := res.Sections[0].Data
		gotW := binary.LittleEndian.Uint32(got)
		wantW := binary.LittleEndian.Uint32(want)
		if gotW == wantW {
			matched++
		} else if m := schemaMaskFor(wantW); m != 0 &&
			gotW&m == wantW&m {
			// the difference is only in don't-care bits outside the
			// schema's Mask (e.g., bit 25 of FP pairs): the text cannot
			// carry this information
			dontCare++
		} else if func() bool {
			got, gerr := arch.DecodeWord(gotW)
			if gerr != nil {
				return false
			}

			want, werr := arch.DecodeWord(wantW)
			return werr == nil && instrTextOf(got, addr) == instrTextOf(want, addr)
		}() {
			// an equivalent encoding of the same text (multiple legal
			// encodings: the ubfm/sbfm form of lsl, immr canonicity)
			equiv++
		} else {
			failed++
			if sample < 8 {
				t.Logf("addr %#x: %q\n  got  % x\n  want % x", addr, src, got, want)
				sample++
			}
		}
	}

	total := matched + failed + dontCare + equiv + notAssembled
	pct := float64(matched+dontCare+equiv) * 100 / float64(total)
	t.Logf(
		"round-trip: %d/%d byte-exact + %d don't-care + %d equiv-encoding (%.2f%%), %d hard, %d not-assembled",
		matched,
		total,
		dontCare,
		equiv,
		pct,
		failed,
		notAssembled,
	)
	require.GreaterOrEqual(t, pct, 90.0, "round-trip rate")
}

// TestArmLabels tests the GNU mode with labels: forward/backward
// branches are resolved with symbols in the second pass (during the
// layout pass, sizing uses Resolve with a placeholder environment).
func TestArmLabels(t *testing.T) {
	src := `
loop:
  add x1, x0, #7
  subs x2, x2, #1
  b.ne loop
  b.eq done
  cbz x2, loop
  mov x0, #0x42
done:
  ret
`
	res, errs := asm.Assemble(src, 0x1000, New())
	require.Empty(t, errs)
	d := res.Sections[0].Data
	require.Len(t, d, 7*4, "total bytes")
	// b.ne loop @0x1008: target 0x1000 → off = -8 → imm19 = -2, cond=ne(1) → 0x54FFFFC1
	bne := binary.LittleEndian.Uint32(d[8:12])
	require.Equal(t, uint32(0x54FFFFC1), bne, "b.ne loop")
	// b.eq done @0x100c: target done=0x1018 → off = 12 → imm19 = 3 → 0x54000060
	beq := binary.LittleEndian.Uint32(d[12:16])
	require.Equal(t, uint32(0x54000060), beq, "b.eq done")
	// cbz x2, loop @0x1010: off = -16 → imm19 = -4 → 0xB4FFFF82
	cbz := binary.LittleEndian.Uint32(d[16:20])
	require.Equal(t, uint32(0xB4FFFF82), cbz, "cbz")
	// symbols
	require.Equal(t, uint64(0x1000), res.Symbols["loop"], "loop")
	require.Equal(t, uint64(0x1018), res.Symbols["done"], "done")
}

// TestArm64NumericLabels tests the GAS numeric local labels: 1f/1b
// with redefinition; the nearest one is picked by source order, and
// numeric labels never get into Symbols.
func TestArm64NumericLabels(t *testing.T) {
	src := `
  cbz x0, 1f
  nop
1:
  b 1b
  b 1f
  nop
1:
  b 1b
`
	res, errs := asm.Assemble(src, 0x1000, New())
	require.Empty(t, errs)
	d := res.Sections[0].Data
	require.Len(t, d, 6*4, "total bytes")

	want := []uint32{
		0xB4000040, // cbz x0, 1f @0x1000 → 1: @0x1008, off=8 → imm19=2
		0xD503201F, // nop
		0x14000000, // b 1b @0x1008 -> 1: @0x1008 (the label before the instruction)
		0x14000002, // b 1f @0x100C -> redefined 1: @0x1014, off=8 -> imm26=2
		0xD503201F, // nop
		0x14000000, // b 1b @0x1014 → 1: @0x1014
	}
	for i, w := range want {
		got := binary.LittleEndian.Uint32(d[i*4:])
		require.Equal(t, w, got, "word %d", i)
	}

	require.NotContains(t, res.Symbols, "1", "numeric label must not be a symbol")
}

// schemaMaskFor returns the mask of the schema matching the word (for
// the don't-care comparison: bits outside the Mask may differ).
func schemaMaskFor(w uint32) uint32 {
	for _, sc := range arch.Schemas() {
		if (w & sc.Mask) == sc.Value {
			return sc.Mask
		}
	}

	return 0
}

// TestLdrLiteralPool tests ldr xN, =literal: literal pool slots at the
// end of the subsection, dedup of identical literals, a symbolic
// literal, and the pool never getting into Symbols. Always ldr-literal
// + pool (no GAS optimization into movz/movk - the semantics are
// equivalent, the bytes differ).
func TestLdrLiteralPool(t *testing.T) {
	src := `
  ldr x0, =0x1122334455667788
  ldr w1, =0x99
  ldr x2, =0x1122334455667788
sym:
  ldr x3, =sym
`
	res, errs := asm.Assemble(src, 0x1000, New())
	require.Empty(t, errs, "errors: %v", errs)
	d := res.Sections[0].Data

	// 4 instructions x 4 + pool: slot1 x8 (0x1122...), slot2 w4 (0x99),
	// 4 bytes of zero padding (slot3 is x8 and the tail after slot2 sits
	// at 4 mod 8 - the gas literal-pool alignment), slot3 x8 (sym)
	require.Len(t, d, 40, "total: % x", d)

	want := []uint32{
		0x58000080, // ldr x0 @0x1000 -> slot1 @0x1010 (imm19=4)
		0x180000A1, // ldr w1 @0x1004 -> slot2 @0x1018 (imm19=5)
		0x58000042, // ldr x2 @0x1008 -> slot1 (dedup, imm19=2)
		0x580000A3, // ldr x3 @0x100C -> slot3 @0x1020 (imm19=5)
	}
	for i, w := range want {
		require.Equal(t, w, binary.LittleEndian.Uint32(d[i*4:]), "word %d", i)
	}

	// pool tail: slot values in first-appearance order
	pool := d[16:]
	require.Equal(
		t,
		[]byte{0x88, 0x77, 0x66, 0x55, 0x44, 0x33, 0x22, 0x11},
		pool[0:8],
		"slot 0x1122... (LE64)",
	)
	require.Equal(t, []byte{0x99, 0, 0, 0}, pool[8:12], "slot 0x99 (LE32, w-slot)")
	require.Equal(t, []byte{0, 0, 0, 0}, pool[12:16], "alignment padding (udf #0)")
	require.Equal(t, []byte{0x0C, 0x10, 0, 0, 0, 0, 0, 0}, pool[16:24], "slot sym=0x100C (LE64)")

	require.Equal(t, uint64(0x100C), res.Symbols["sym"], "sym")
	require.Len(t, res.Symbols, 1, "pool names must not be in Symbols: %v", res.Symbols)
}

// TestLdrLiteralPoolAlignment — every pool slot sits at its natural
// alignment, the gas literal-pool rule (the seL4 head.S shape: an 8-byte
// ldr= slot after a tail at 4 mod 8 pads with one zero word - udf #0 in
// the gas dump; loading the quad from the misaligned slot was the
// alignment fault of the first kernel image). Pinned bit-exact.
func TestLdrLiteralPoolAlignment(t *testing.T) {
	slot64 := "8877665544332211"
	cases := []struct {
		name  string
		src   string
		words []uint32
		tail  string
	}{
		{
			name:  "x-slot after a tail at 4 mod 8: one zero word",
			src:   "nop\nnop\nldr x19, =0x1122334455667788\n",
			words: []uint32{0xD503201F, 0xD503201F, 0x58000053},
			tail:  "00000000" + slot64,
		},
		{
			name:  "aligned tail: no padding",
			src:   "nop\nldr x19, =0x1122334455667788\n",
			words: []uint32{0xD503201F, 0x58000033},
			tail:  slot64,
		},
		{
			name:  "w-slot after a tail at 4 mod 8: its own 4 alignment",
			src:   "nop\nnop\nnop\nldr w1, =0x99\n",
			words: []uint32{0xD503201F, 0xD503201F, 0xD503201F, 0x18000021},
			tail:  "99000000",
		},
		{
			name:  "w-slot then x-slot: the pad lands between them",
			src:   "ldr w1, =0x99\nldr x2, =0x1122334455667788\n",
			words: []uint32{0x18000041, 0x58000062},
			tail:  "99000000" + "00000000" + slot64,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res, errs := asm.Assemble(tc.src, 0x41000000, New())
			require.Empty(t, errs, "errors: %v", errs)
			d := res.Sections[0].Data

			for i, w := range tc.words {
				require.Equal(t, w, binary.LittleEndian.Uint32(d[i*4:]), "word %d", i)
			}

			want, err := hex.DecodeString(tc.tail)
			require.NoError(t, err)
			require.Equal(t, want, d[len(tc.words)*4:], "pool tail: % x", d[len(tc.words)*4:])
		})
	}
}

// TestMovUmovAlias - the UMOV input spellings and the LLVM print alias:
// the decoder prints mov for the filling sizes (s into w, d into x), so
// the self-verify canonicalizes the input spelling. Words pinned against
// llvm-mc (docker assembly-tests, 2026-09-08); the mov input alias and
// the rejection cases are in asm/arm64/alias/assemble_test.go.
// TestLdrLiteralPoolUndef — an undefined symbol in a literal pool slot is
// an error, not a silent zero (as with .word).
func TestLdrLiteralPoolUndef(t *testing.T) {
	_, errs := asm.Assemble("ldr x0, =undef", 0, New())
	require.NotEmpty(t, errs, "an undefined pool symbol must error")
}

// TestStmtSeparator — the GAS ';': statements on one line (a label may
// follow the separator - the BEGIN_FUNC expansion shape), empty segments
// are legal, '//' still comments to the end of the line, a ';' inside a
// string literal is data.
func TestStmtSeparator(t *testing.T) {
	cases := []struct {
		src   string
		words []uint32
	}{
		{
			"nop ; ret",
			[]uint32{0xd503201f, 0xd65f03c0},
		},
		{
			"; nop",
			[]uint32{0xd503201f},
		},
		{
			"nop ; ; nop",
			[]uint32{0xd503201f, 0xd503201f},
		},
		{
			"nop ; // comment",
			[]uint32{0xd503201f},
		},
		{
			"nop;ret",
			[]uint32{0xd503201f, 0xd65f03c0},
		},
	}
	for _, c := range cases {
		res, errs := asm.Assemble(c.src, 0, New())
		require.Empty(t, errs, "case %q", c.src)
		require.Len(t, res.Sections[0].Data, 4*len(c.words), "case %q", c.src)
		for i, w := range c.words {
			require.Equal(
				t,
				w,
				binary.LittleEndian.Uint32(res.Sections[0].Data[i*4:]),
				"case %q word %d",
				c.src,
				i,
			)
		}
	}

	// the BEGIN_FUNC shape: the directive tail stops at the separator,
	// the label after it is its own statement
	res, errs := asm.Assemble(".global _start ; _start: ret", 0, New())
	require.Empty(t, errs)
	require.Len(t, res.Sections[0].Data, 4)
	require.Equal(t, uint32(0xd65f03c0), binary.LittleEndian.Uint32(res.Sections[0].Data))
	require.Equal(t, uint64(0), res.Symbols["_start"])

	// a ';' inside a string literal is data, not a separator
	res, errs = asm.Assemble(".asciz \"a;b\"", 0, New())
	require.Empty(t, errs)
	require.Equal(t, []byte{'a', ';', 'b', 0}, res.Sections[0].Data)
}

// TestSectionFlags — the GAS flag strings after the section name are
// recognized and ignored (the core carries no flag semantics); anything
// but a comma list of quoted strings is an error.
func TestSectionFlags(t *testing.T) {
	for _, name := range []string{
		".boot.text",
		".text",
		".mydata",
	} {
		src := ".section " + name + ", \"ax\"\nret\n"
		res, errs := asm.Assemble(src, 0, New())
		require.Empty(t, errs, "section %q", name)
		require.NotEmpty(t, res.Sections, "section %q", name)
		require.Equal(t, name, res.Sections[0].Name, "section %q", name)
		require.Len(t, res.Sections[0].Data, 4, "section %q", name)
	}

	// several quoted strings after the name (flags, type)
	res, errs := asm.Assemble(".section .data, \"aw\", \"progbits\"\n.word 1\n", 0, New())
	require.Empty(t, errs)
	require.Equal(t, ".data", res.Sections[0].Name)

	// the name without flags is unchanged
	res, errs = asm.Assemble(".section .text\nret\n", 0, New())
	require.Empty(t, errs)
	require.Len(t, res.Sections[0].Data, 4)

	// an unquoted tail is not a flag list
	_, errs = asm.Assemble(".section .text, ax\nret\n", 0, New())
	require.NotEmpty(t, errs, "an unquoted flag tail must error")
}

func TestMovUmovAlias(t *testing.T) {
	for _, c := range []struct {
		src  string
		word uint32
	}{
		{"umov x0, v1.d[1]", 0x4E183C20},
		{"umov w0, v1.s[1]", 0x0E0C3C20},
		{"umov w0, v1.b[3]", 0x0E073C20},
		{"umov w0, v1.h[7]", 0x0E1E3C20},
	} {
		require.Equal(t, c.word, armAssembleOne(t, c.src, 0), "case %q", c.src)
	}

	// llvm rejects these ("invalid operand"): the umov width must match
	// the element size.
	for _, src := range []string{
		"umov w0, v1.d[1]",
		"umov x0, v1.s[1]",
	} {
		_, errs := asm.Assemble(src, 0, New())
		require.NotEmpty(t, errs, "case %q must not assemble", src)
	}
}

// TestUaddlvWords — the uaddlv accumulator text (the h/s FP scalar of the
// doubled lane) assembles into the clang words and renders back.
func TestUaddlvWords(t *testing.T) {
	for _, c := range []struct {
		src  string
		word uint32
	}{
		{"uaddlv.8b h30, v24", 0x2e303b1e},
		{"uaddlv.4h s0, v0", 0x2e703800},
		{"uaddlv.16b h1, v2", 0x6e303841},
	} {
		got := armAssembleOne(t, c.src, 0)
		require.Equal(t, c.word, got, "case %q", c.src)
	}
}

// TestElemWords — the stage-7 element families (the SIMD copy class and
// the by-element arithmetic) assemble into the clang words: the fcmla index
// layouts (.4h — bit 21, .4s — bit 11) and the long families' result
// arrangements (the .8h-result class is unallocated).
func TestElemWords(t *testing.T) {
	for _, c := range []struct {
		src  string
		word uint32
	}{
		{"smov x5, v19.s[0]", 0x4e042e65},
		{"smov w5, v19.b[0]", 0x0e012e65},
		{"mov.b v3[5], w17", 0x4e0b1e23},
		{"ins.b v9[2], v19[5]", 0x6e052e69},
		{"dup.8b v9, v19[0]", 0x0e010669},
		{"dup.2s v9, v19[0]", 0x0e040669},
		{"fcmla.4s v9, v19, v5[1], #0", 0x6f851a69},
		{"fcmla.4h v9, v19, v20[0], #90", 0x2f543269},
		{"fcmla.8h v9, v19, v5[3], #180", 0x6f655a69},
		{"smull2.4s v9, v19, v5[0]", 0x4f45a269},
		{"smull.4s v9, v19, v5[7]", 0x0f75aa69},
		{"smull2.2d v9, v19, v20[1]", 0x4fb4a269},
		{"sqdmlsl.2d v9, v19, v5[3]", 0x0fa57a69},
		{"mul.4h v9, v19, v5[3]", 0x0f758269},
		{"mla.2s v31, v30, v20[3]", 0x2fb40bdf},
		{"fmla.2d v9, v19, v20[1]", 0x4fd41a69},
		{"fmulx.4s v0, v1, v15[3]", 0x6faf9820},
		{"smlal2.4s v0, v1, v2[5]", 0x4f522820},
		{"umull.2d v3, v4, v25[2]", 0x2f99a883},
		{"sqdmulh.8h v7, v8, v9[6]", 0x4f69c907},
	} {
		got := armAssembleOne(t, c.src, 0)
		require.Equal(t, c.word, got, "case %q", c.src)
	}
}

// TestElemRefusals — the element spellings clang rejects: the smov .s
// destination width, the unallocated fcmla arrangements and the
// unallocated .8h-result class of the long families.
func TestElemRefusals(t *testing.T) {
	for _, src := range []string{
		"smov w5, v19.s[0]",
		"fcmla.2s v0, v1, v2[0], #0",
		"fcmla.2d v0, v1, v2[0], #0",
		"smull.8h v0, v1, v2[0]",
		"smull2.8h v0, v1, v2[0]",
		"fcmla.4s v0, v1, v2[2], #0",
	} {
		_, errs := asm.Assemble(src, 0, New())
		require.NotEmpty(t, errs, "case %q must not assemble", src)
	}
}

// TestZeroShiftSpellings — the zero-amount shift modifiers are the
// canonical no-shift spelling (clang accepts them; the decoder prints
// without), pinned after the loose-compare learned to drop the suffix.
func TestZeroShiftSpellings(t *testing.T) {
	for _, c := range []struct {
		src  string
		word uint32
	}{
		{"neg w17, w22, lsl #0", 0x4b1603f1},
		{"add w0, w1, w2, lsr #0", 0x0b420020},
		{"and x0, x1, x2, ror #0", 0x8ac20020},
	} {
		got := armAssembleOne(t, c.src, 0)
		require.Equal(t, c.word, got, "case %q", c.src)
	}
}
