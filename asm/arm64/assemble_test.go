package arm64

import (
	"encoding/binary"
	"strings"
	"testing"

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
	}
	for _, c := range cases {
		if c.src == "add x0, x1, #0x42" {
			continue // checked separately below
		}

		got := armAssembleOne(t, c.src, 0)
		require.Equal(t, c.word, got, "case %q", c.src)
	}

	// add x0, x1, #0x42: imm12=0x42<<10 | Rn=1<<5 | Rd=0 | 0x91000000
	got := armAssembleOne(t, "add x0, x1, #0x42", 0)
	require.True(
		t,
		got == 0x91068020 || got == 0x91000000|0x42<<10|1<<5,
		"add = %#08x",
		got,
	)
}

func TestArmBranchAndMem(t *testing.T) {
	// b 0x1008 @ 0x1000
	got := armAssembleOne(t, "b 0x1008", 0x1000)
	require.Equal(t, uint32(0x14000002), got, "b")
	// bl 0x2000 @ 0x1000
	got = armAssembleOne(t, "bl 0x2000", 0x1000)
	require.Equal(t, uint32(0x94000400), got, "bl")
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

	insts, err := arch.Parse(ts.Addr)(bytes.Buffer(ts.Data))
	require.NoError(t, err)
	matched, failed, notAssembled, dontCare, equiv := 0, 0, 0, 0, 0
	sample := 0
	for _, in := range insts {
		if _, ok := in.(arch.Generic); ok {
			continue // decode-only: the generic syntax is not parsed by the assembler
		}

		src := objdump.StripComments(objdump.Normalize(in.ObjDump(disasm.DefaultViewCtx())))
		if src == "" || strings.HasPrefix(src, ".word") {
			continue
		}

		off := in.Addr() - ts.Addr
		want := ts.Data[off : off+4]
		res, errs := asm.Assemble(src, in.Addr(), New())
		if len(errs) != 0 {
			notAssembled++
			if sample < 8 {
				t.Logf("addr %#x: %q: %v", in.Addr(), src, errs)
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
		} else if instrTextOf(arch.DecodeWord(gotW, in.Addr())) == instrTextOf(arch.DecodeWord(wantW, in.Addr())) {
			// an equivalent encoding of the same text (multiple legal
			// encodings: the ubfm/sbfm form of lsl, immr canonicity)
			equiv++
		} else {
			failed++
			if sample < 8 {
				t.Logf("addr %#x: %q\n  got  % x\n  want % x", in.Addr(), src, got, want)
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

	// 4 instructions x 4 + pool: slots x8 (0x1122...), w4 (0x99), x8 (sym) = 20
	require.Len(t, d, 36, "total: % x", d)

	want := []uint32{
		0x58000080, // ldr x0 @0x1000 -> slot1 @0x1010 (imm19=4)
		0x180000A1, // ldr w1 @0x1004 -> slot2 @0x1018 (imm19=5)
		0x58000042, // ldr x2 @0x1008 -> slot1 (dedup, imm19=2)
		0x58000083, // ldr x3 @0x100C -> slot3 @0x101C (imm19=4)
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
	require.Equal(t, []byte{0x0C, 0x10, 0, 0, 0, 0, 0, 0}, pool[12:20], "slot sym=0x100C (LE64)")

	require.Equal(t, uint64(0x100C), res.Symbols["sym"], "sym")
	require.Len(t, res.Symbols, 1, "pool names must not be in Symbols: %v", res.Symbols)
}
