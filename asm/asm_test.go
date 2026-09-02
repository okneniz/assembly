package asm

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/okneniz/parsec"
	parsecstrings "github.com/okneniz/parsec/strings"
	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/asm/expr"
)

var errMock = errors.New("mock encode failure")

// mockInstr / failInstr - mock unresolved instructions: "pad N" resolves to
// N bytes of 0xAB, "fail" fails at resolution.
type mockInstr struct {
	n int
}

func newMockInstr(n int) mockInstr {
	return mockInstr{n: n}
}

func (m mockInstr) Resolve(Ctx) (Resolved, error) {
	return mockResolved{b: bytes.Repeat([]byte{0xAB}, m.n)}, nil
}

type mockResolved struct {
	b []byte
}

func (m mockResolved) Encode(w io.Writer) (int64, error) {
	n, err := w.Write(m.b)
	return int64(n), err
}

type failInstr struct{}

func (failInstr) Resolve(Ctx) (Resolved, error) {
	return nil, errMock
}

// mockBackend - a minimal Syntax for core tests.
type mockBackend struct{}

func (mockBackend) Instruction() parsec.Combinator[rune, parsecstrings.Position, Unresolved] {
	pad := parsecstrings.Cast(
		parsecstrings.Skip(
			parsecstrings.String("pad", "pad"),
			parsecstrings.SkipMany(expr.CSpace,
				parsecstrings.Cast(
					parsecstrings.Some(4, "operand", expr.CDecDigit()),
					func(rs []rune) (string, error) {
						return string(rs), nil
					},
				)),
		),
		func(s string) (Unresolved, error) {
			n := 0
			for _, r := range s {
				n = n*10 + int(r-'0')
			}

			return newMockInstr(n), nil
		},
	)
	fail := parsecstrings.Cast(
		parsecstrings.String("fail", "fail"),
		func(string) (Unresolved, error) {
			return failInstr{}, nil
		},
	)
	pool := parsecstrings.Cast(
		parsecstrings.Skip(
			parsecstrings.String("pool", "pool"),
			parsecstrings.SkipMany(expr.CSpace,
				parsecstrings.Cast(
					parsecstrings.Some(4, "operand", expr.CDecDigit()),
					func(rs []rune) (string, error) {
						return string(rs), nil
					},
				)),
		),
		func(s string) (Unresolved, error) {
			n := int64(0)
			for _, r := range s {
				n = n*10 + int64(r-'0')
			}

			return mockPoolInstr{v: expr.Num(n * 0x101)}, nil
		},
	)

	return parsecstrings.Choice(
		"instruction",
		parsecstrings.Try(pad),
		parsecstrings.Try(pool),
		parsecstrings.Try(fail),
	)
}

// mockPoolInstr - the "pool N" instruction for literal pool unit tests:
// it requires a slot with the value N*0x101 (PoolUser), encodes 4 bytes -
// the low word of the ADDRESS OF ITS OWN SLOT (the reserved name PoolSelf).
type mockPoolInstr struct {
	v *expr.Expr
}

func (m mockPoolInstr) Resolve(c Ctx) (Resolved, error) {
	addr, ok := c.Resolve(PoolSelf)
	if !ok {
		return nil, errors.New("pool slot not resolved")
	}

	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], uint32(addr))
	return mockResolved{b: buf[:]}, nil
}

func (m mockPoolInstr) PoolReq() (*expr.Expr, int, bool) {
	return m.v, 8, true
}

func (mockBackend) Comment() parsec.Combinator[rune, parsecstrings.Position, string] {
	hash := parsecstrings.Cast(
		parsecstrings.Skip(
			parsecstrings.Try(parsecstrings.Eq("comment", '#')),
			parsecstrings.Many(4, cNotNL),
		),
		func(rs []rune) (string, error) {
			return string(rs), nil
		},
	)
	slash := parsecstrings.Cast(
		parsecstrings.Skip(
			parsecstrings.Try(parsecstrings.String("comment", "//")),
			parsecstrings.Many(4, cNotNL),
		),
		func(rs []rune) (string, error) {
			return string(rs), nil
		},
	)
	// Choice alternatives must be Try-wrapped (the parsec contract)
	return parsecstrings.Choice("comment", parsecstrings.Try(slash), parsecstrings.Try(hash))
}

// mockOpts - the applied .option values (for ApplyOption tests).
var mockOpts []string

func (mockBackend) ApplyOption(name string) error {
	mockOpts = append(mockOpts, name)
	return nil
}

func (mockBackend) ResetOptions() {
	mockOpts = nil
}

func TestSymbolsAndSections(t *testing.T) {
	src := `
start:
  pad 4
loop:
  pad 8
  .data
val: .word 0x12345678, start
  .text
  pad 2
`
	res, errs := Assemble(src, 0x1000, mockBackend{})
	require.Empty(t, errs, "unexpected errors: %v", errs)
	require.Len(t, res.Sections, 2, "sections: %+v", res.Sections)
	text, data := res.Sections[0], res.Sections[1]
	require.Equal(t, uint64(0x1000), text.Addr, "text addr")
	require.Len(t, text.Data, 14, "text len")
	require.Equal(t, uint64(0x1000+14), data.Addr, "data addr")
	require.Len(t, data.Data, 8, "data len")
	require.Equal(t, []byte{0x78, 0x56, 0x34, 0x12}, data.Data[0:4], ".word LE")
	// start = 0x1000 - the base address of the first section
	start := uint64(
		data.Data[4],
	) | uint64(
		data.Data[5],
	)<<8 | uint64(
		data.Data[6],
	)<<16 | uint64(
		data.Data[7],
	)<<24
	require.Equal(t, uint64(0x1000), start, ".word start")
	require.Equal(t, uint64(0x1004), res.Symbols["loop"], "loop")
}

func TestDirectivesLayout(t *testing.T) {
	src := `
  .byte 1, 2
  .half 0x0102
  .word -1
  .quad 0x1122334455667788
  .string "ab", "c"
  .ascii "xy"
  .zero 3
  .balign 4
`
	res, errs := Assemble(src, 0, mockBackend{})
	require.Empty(t, errs, "errors: %v", errs)
	require.NotEmpty(t, res.Sections, "no sections assembled")
	d := res.Sections[0].Data
	// 2 + 2 + 4 + 8 + 5 + 2 + 3 = 26; .balign 4 pads up to 28
	require.Len(t, d, 28, "len: % x", d)
	require.Equal(t, []byte{1, 2, 0x02, 0x01}, d[:4], "byte/half")
	require.Equal(t, []byte{0xff, 0xff, 0xff, 0xff}, d[4:8], ".word -1")
	require.Equal(t, byte(0x88), d[8], ".quad LE low")
	require.Equal(t, byte(0x11), d[15], ".quad LE high")
	require.Equal(t, "ab\x00c\x00", string(d[16:21]), ".string")
	require.Equal(t, "xy", string(d[21:23]), ".ascii")
	for _, b := range d[23:] {
		require.Zero(t, b, "padding not zero: % x", d[23:])
	}
}

func TestAlign(t *testing.T) {
	res, errs := Assemble(".byte 1\n.align 3\n.byte 2\n", 0, mockBackend{})
	require.Empty(t, errs, "errors: %v", errs)
	// .align 3 = a power of two: 8; after 1 byte - 7 zeros
	require.Len(t, res.Sections[0].Data, 9)

	// .p2align - the same power-of-two alignment; a name with a digit
	// parses (cIdent)
	res, errs = Assemble(".byte 1\n.p2align 3\n.byte 2\n", 0, mockBackend{})
	require.Empty(t, errs, "errors: %v", errs)
	require.Len(t, res.Sections[0].Data, 9)
}

func TestAlignBoundaries(t *testing.T) {
	// extreme values must not panic (division by zero) and must not
	// duplicate the error between passes
	cases := []struct {
		src     string
		wantMsg string
	}{
		{".align 64\n", "alignment 2^64 too large"},   // 1<<64 = 0 in Go
		{".p2align 64\n", "alignment 2^64 too large"}, // the same for p2align
		{".balign 0\n", "zero alignment"},             // degenerate alignment
		{".align -1\n", "constant alignment required"},
	}

	for _, c := range cases {
		_, errs := Assemble(c.src, 0, mockBackend{})
		require.Len(t, errs, 1, "Assemble(%q) errs = %v", c.src, errs)
		require.Contains(t, errs[0].Msg, c.wantMsg, "Assemble(%q)", c.src)
	}

	// .align 0 = 2^0 = 1 - legal, no padding
	res, errs := Assemble(".byte 1\n.align 0\n.byte 2\n", 0, mockBackend{})
	require.Empty(t, errs, "errors: %v", errs)
	require.Len(t, res.Sections[0].Data, 2)
}

func TestSetAndDot(t *testing.T) {
	src := `
  .set base, 0x40
target: .word base + 4, . + 8
`
	res, errs := Assemble(src, 0x100, mockBackend{})
	require.Empty(t, errs, "errors: %v", errs)
	d := res.Sections[0].Data
	w0 := le32(d[0:4])
	w1 := le32(d[4:8])
	require.Equal(t, uint64(0x44), w0, "base+4")
	// "." = the address of .word (0x100), .+8 = 0x108
	require.Equal(t, uint64(0x100+8), w1, ". + 8")
}

func TestErrorPositionsAndRecovery(t *testing.T) {
	src := "pad 4\nbadsyntax here\npad 2\n"
	res, errs := Assemble(src, 0, mockBackend{})
	require.Len(t, errs, 1, "errs = %v", errs)
	require.Equal(t, uint(2), errs[0].Line, "error line (%v)", errs[0])
	require.Len(t, res.Sections[0].Data, 6, "recovered len")
}

func TestUnknownDirective(t *testing.T) {
	_, errs := Assemble(".frobnicate 1\n", 0, mockBackend{})
	require.Len(t, errs, 1, "errs = %v", errs)
	require.Contains(t, errs[0].Msg, "unknown directive")
}

func TestCommentStyles(t *testing.T) {
	src := "# full line comment\npad 4 # trailing\n// slash comment\npad 2 // trailing slash\n"
	res, errs := Assemble(src, 0, mockBackend{})
	require.Empty(t, errs, "errors: %v", errs)
	require.Len(t, res.Sections[0].Data, 6)
}

func TestIgnoredDirectives(t *testing.T) {
	src := ".file 1 \"a.c\"\n.loc 1 2 0\n.cfi_startproc\npad 4\n.cfi_endproc\n"
	res, errs := Assemble(src, 0, mockBackend{})
	require.Empty(t, errs, "errors: %v", errs)
	require.Len(t, res.Sections[0].Data, 4)
}

func TestLabelRedefinition(t *testing.T) {
	_, errs := Assemble("x: pad 2\nx: pad 2\n", 0, mockBackend{})
	require.Len(t, errs, 1, "errs = %v", errs)
	require.Contains(t, errs[0].Msg, "redefined")
}

func TestNumericLabels(t *testing.T) {
	cases := []struct {
		name     string
		src      string
		base     uint64
		wantData []byte // the exact image of .text
	}{
		{
			"f forward, b backward, redefinition of the same label",
			`
0:
  pad 4
  .word 0f
0:
  .word 0b
  .word 0f
0:
`,
			0x100,
			[]byte{
				0xAB, 0xAB, 0xAB, 0xAB, // pad 4
				0x08, 0x01, 0x00, 0x00, // 0f → the next 0: @0x108
				0x08, 0x01, 0x00, 0x00, // 0b → 0: @0x108
				0x10, 0x01, 0x00, 0x00, // 0f → the redefined 0: @0x110
			},
		},
		{
			"b on the definition line - the label before the instruction",
			`1: .word 1b`,
			0x100,
			[]byte{0x00, 0x01, 0x00, 0x00},
		},
		{
			"binary literal and reference in one list",
			`0: .word 0b1010, 0b`,
			0x100,
			[]byte{
				0x0A, 0x00, 0x00, 0x00, // 0b1010 - literal
				0x00, 0x01, 0x00, 0x00, // 0b - reference to 0: @0x100
			},
		},
		{
			"multi-digit numbers and different labels mixed",
			`
10:
  .word 10f, 2f
2: pad 2
10: .word 2b
`,
			0x100,
			[]byte{
				0x0A, 0x01, 0x00, 0x00, // 10f → the second 10: @0x10A
				0x08, 0x01, 0x00, 0x00, // 2f → 2: @0x108
				0xAB, 0xAB, // pad 2
				0x08, 0x01, 0x00, 0x00, // 2b → 2: @0x108
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res, errs := Assemble(c.src, c.base, mockBackend{})
			require.Empty(t, errs, "unexpected errors: %v", errs)
			require.NotEmpty(t, res.Sections, "no sections assembled")
			require.Equal(t, c.wantData, res.Sections[0].Data, "data: % x", res.Sections[0].Data)

			for _, name := range []string{"0", "1", "2", "10"} {
				require.NotContains(
					t,
					res.Symbols,
					name,
					"numeric label %q must not be a symbol",
					name,
				)
			}
		})
	}
}

func TestNumericLabelUndefined(t *testing.T) {
	cases := []struct {
		src     string
		wantMsg string
	}{
		{`  .word 5b`, `undefined symbol "5b"`},  // no definition at all
		{`1: .word 1f`, `undefined symbol "1f"`}, // no forward definitions left
		{`0: .word 1b`, `undefined symbol "1b"`}, // there is 0:, no 1:
	}

	for _, c := range cases {
		_, errs := Assemble(c.src, 0, mockBackend{})
		require.Len(t, errs, 1, "Assemble(%q) errs = %v", c.src, errs)
		require.Contains(t, errs[0].Msg, c.wantMsg, "Assemble(%q)", c.src)
	}
}

func TestEncodeErrorFillsZeros(t *testing.T) {
	_, errs := Assemble("fail\npad 2\n", 0, mockBackend{})
	require.NotEmpty(t, errs, "want encode error")
}

func TestSkipDirective(t *testing.T) {
	// .skip - an alias of .zero/.space: N zero bytes
	res, errs := Assemble(".byte 1\n.skip 3\n.byte 2\n", 0, mockBackend{})
	require.Empty(t, errs, "errors: %v", errs)
	require.Equal(t, []byte{1, 0, 0, 0, 2}, res.Sections[0].Data)
}

func TestEndDirective(t *testing.T) {
	// .end - the end of the source: lines below are ignored entirely,
	// including parse errors and encoding
	res, errs := Assemble("pad 2\n.end\npad 2\n@@junk\nfail\n", 0, mockBackend{})
	require.Empty(t, errs, "post-.end lines: %v", errs)
	require.Len(t, res.Sections[0].Data, 2)
}

func TestOptionDirective(t *testing.T) {
	// .option - a semantic directive: the values go to Syntax.ApplyOption
	// (in order of appearance), the reset between passes does not duplicate
	// the log
	_, errs := Assemble(".option push\n.option norvc\npad 2\n.option pop\n", 0, mockBackend{})
	require.Empty(t, errs, "errors: %v", errs)
	require.Equal(t, []string{"push", "norvc", "pop"}, mockOpts)
}

func TestBss(t *testing.T) {
	src := `
  pad 4
  .bss
buf: .zero 100
  .align 5
  .skip 7
end:
`
	res, errs := Assemble(src, 0x1000, mockBackend{})
	require.Empty(t, errs, "errors: %v", errs)
	require.Len(t, res.Sections, 2, "sections: %+v", res.Sections)

	text, bss := res.Sections[0], res.Sections[1]
	require.Equal(t, ".text", text.Name)
	require.Len(t, text.Data, 4)
	require.False(t, text.Nobits)

	require.Equal(t, ".bss", bss.Name)
	require.True(t, bss.Nobits, "bss is NOBITS")
	require.Nil(t, bss.Data, "NOBITS has no file data")
	// 100 + padding up to 32 (28) + 7
	require.Equal(t, 135, bss.Size, "bss mem size")
	require.Equal(t, uint64(0x1004), res.Symbols["buf"], "buf")
	require.Equal(t, uint64(0x1004+135), res.Symbols["end"], "end")
}

func TestBssRejectsData(t *testing.T) {
	cases := []struct {
		src     string
		wantMsg string
	}{
		{".bss\nx: .word 5\n", "only .zero/.skip/.align permitted"},
		{".bss\nx: .byte 1\n", "only .zero/.skip/.align permitted"},
		{".bss\nx: .string \"a\"\n", "only .zero/.skip/.align permitted"},
		{".bss\npad 2\n", "instructions are not permitted"},
	}

	for _, c := range cases {
		_, errs := Assemble(c.src, 0, mockBackend{})
		require.Len(t, errs, 1, "Assemble(%q) errs = %v", c.src, errs)
		require.Contains(t, errs[0].Msg, c.wantMsg, "Assemble(%q)", c.src)
	}
}

func TestIncbin(t *testing.T) {
	path := filepath.Join(t.TempDir(), "blob.bin")
	require.NoError(t, os.WriteFile(path, []byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE}, 0o644))

	cases := []struct {
		src      string
		wantData []byte
	}{
		{".byte 1\n.incbin \"" + path + "\"\n", []byte{1, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE}},
		{".incbin \"" + path + "\", 2\n", []byte{0xCC, 0xDD, 0xEE}}, // skip
		{".incbin \"" + path + "\", 1, 2\n", []byte{0xBB, 0xCC}},    // skip+count
		{
			".incbin \"" + path + "\", 4, 100\n",
			[]byte{0xEE},
		}, // count past the end - truncated
		{".incbin \"" + path + "\", 100\n", nil},  // skip past the end - empty
		{".incbin \"" + path + "\", 0, 0\n", nil}, // zero length
	}

	for _, c := range cases {
		res, errs := Assemble(c.src, 0, mockBackend{})
		require.Empty(t, errs, "Assemble(%q): %v", c.src, errs)
		if c.wantData == nil {
			// an empty insertion - the section drops out (the empty filter)
			for _, s := range res.Sections {
				require.Empty(t, s.Data, "Assemble(%q)", c.src)
			}

			continue
		}

		require.Equal(t, c.wantData, res.Sections[0].Data, "Assemble(%q)", c.src)
	}
}

func TestIncbinErrors(t *testing.T) {
	_, errs := Assemble(".incbin \"/nonexistent/blob.bin\"\n", 0, mockBackend{})
	require.Len(t, errs, 1, "errs = %v", errs)
	require.Contains(t, errs[0].Msg, ".incbin:", "no file - error")

	_, errs = Assemble(".bss\n.incbin \"whatever\"\n", 0, mockBackend{})
	require.Len(t, errs, 1, "errs = %v", errs)
	require.Contains(t, errs[0].Msg, "NOBITS", "bss - error")
}

func le32(b []byte) uint64 {
	return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24
}

func TestSubsections(t *testing.T) {
	src := `
.text 1
hi: pad 4
.text 0
lo: pad 2
.text 1
  pad 2
tail:
`
	res, errs := Assemble(src, 0x1000, mockBackend{})
	require.Empty(t, errs, "errors: %v", errs)
	require.Len(t, res.Sections, 1, "sections: %+v", res.Sections)

	// GAS concatenation: subsection 0 (2 bytes), then 1 (4+2); the data is
	// the concatenation
	require.Equal(
		t,
		[]byte{0xAB, 0xAB, 0xAB, 0xAB, 0xAB, 0xAB, 0xAB, 0xAB},
		res.Sections[0].Data,
		"data / sub0(2)+sub1(6)",
	)
	require.Equal(t, uint64(0x1000), res.Symbols["lo"], "lo @ sub0")
	require.Equal(t, uint64(0x1002), res.Symbols["hi"], "hi @ sub1")
	require.Equal(t, uint64(0x1008), res.Symbols["tail"], "tail @ end of sub1")
}

func TestSubsectionsDataAndBss(t *testing.T) {
	src := `
.data 1
d1: .word 0x11
.data 0
d0: .word 0x22
.bss 2
b2: .zero 8
.bss
b0: .zero 4
`
	res, errs := Assemble(src, 0x1000, mockBackend{})
	require.Empty(t, errs, "errors: %v", errs)
	require.Len(t, res.Sections, 2, "sections: %+v", res.Sections)

	data, bss := res.Sections[0], res.Sections[1]
	require.Equal(t, ".data", data.Name)
	require.Equal(
		t,
		[]byte{0x22, 0, 0, 0, 0x11, 0, 0, 0},
		data.Data,
		".data concatenation by numbers",
	)
	require.Equal(t, uint64(0x1000), res.Symbols["d0"])
	require.Equal(t, uint64(0x1004), res.Symbols["d1"])

	require.Equal(t, ".bss", bss.Name)
	require.True(t, bss.Nobits)
	require.Equal(t, 12, bss.Size, ".bss subsections concatenated")
	require.Equal(t, uint64(0x1008), res.Symbols["b0"], ".bss after .data")
	require.Equal(t, uint64(0x100C), res.Symbols["b2"], "bss sub2 after sub0")
}

func TestSubsectionsErrors(t *testing.T) {
	cases := []struct {
		src     string
		wantMsg string
	}{
		{".text 8193\n", "subsection number must be a constant 0..8192"},
		{".data -1\n", "subsection number must be a constant 0..8192"},
		{".bss undefined_sym\n", "subsection number must be a constant 0..8192"},
	}

	for _, c := range cases {
		_, errs := Assemble(c.src, 0, mockBackend{})
		require.Len(t, errs, 1, "Assemble(%q) errs = %v", c.src, errs)
		require.Contains(t, errs[0].Msg, c.wantMsg, "Assemble(%q)", c.src)
	}
}

func TestSubsectionsNumericLabels(t *testing.T) {
	// numeric labels across subsections: 1: in sub1 and 1: in sub0;
	// references from sub0 - both to the nearest in SOURCE ORDER (the sub1
	// definition is written earlier), the address is final per the GAS
	// concatenation: sub1 starts after sub0 (8 bytes)
	src := `
.text 1
1: pad 4
.text 0
  .word 1f
  .word 1b
1: pad 2
`
	res, errs := Assemble(src, 0x1000, mockBackend{})
	require.Empty(t, errs, "errors: %v", errs)
	d := res.Sections[0].Data
	// sub0 = two words + pad 2 (label 1: of line 7) = 10 bytes; sub1
	// (pad 4, label 1: of line 3) after concatenation x100A.
	// 1f from line 5 - forward: the label of line 7 (sub0) x1008;
	// 1b - backward: the label of line 3 (sub1) x100A
	require.Equal(t, uint64(0x1008), le32(d[0:4]), "1f")
	require.Equal(t, uint64(0x100A), le32(d[4:8]), "1b")
}

func TestLiteralPool(t *testing.T) {
	// the core literal pool mechanics (mock backend): slots go at the tail
	// of the subsection, dedup by value, slot addresses resolve via the
	// auto-name, the pool does not get into Symbols
	src := "pool 1\npool 2\npool 1\n"
	res, errs := Assemble(src, 0x1000, mockBackend{})
	require.Empty(t, errs, "errors: %v", errs)
	require.Len(t, res.Sections, 1)

	d := res.Sections[0].Data
	// 3 instructions × 4 bytes + the pool: slots 0x101 and 0x202 (dedup)
	// × 8 bytes
	require.Len(t, d, 12+16, "data: % x", d)

	// the instructions encoded the ADDRESSES of their slots (LE32):
	// slot 0x101 @0x100C, slot 0x202 @0x1014, the repeat of 0x101 @0x100C
	require.Equal(t, uint64(0x100C), le32(d[0:4]), "pool 1 → slot 1")
	require.Equal(t, uint64(0x1014), le32(d[4:8]), "pool 2 → slot 2")
	require.Equal(t, uint64(0x100C), le32(d[8:12]), "dedup: the same slot")

	// the pool tail: slot values LE64 in first-appearance order
	require.Equal(t, []byte{0x01, 0x01, 0, 0, 0, 0, 0, 0}, d[12:20], "slot 0x101")
	require.Equal(t, []byte{0x02, 0x02, 0, 0, 0, 0, 0, 0}, d[20:28], "slot 0x202")

	require.NotContains(t, res.Symbols, poolName(8, "257"), "the pool is not in Symbols")
}

func TestLtorgUnsupported(t *testing.T) {
	_, errs := Assemble(".ltorg\n", 0, mockBackend{})
	require.Len(t, errs, 1, "errs = %v", errs)
	require.Contains(t, errs[0].Msg, "end of the subsection", ".ltorg - an explicit error")
}
