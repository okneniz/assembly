package assembly_test

// Property tests (oh-snap), loong64: round trip of single instructions
// and compositions (fixed 32-bit words, canonicalization of the alias
// forms), decoder robustness, differential against the real objdump,
// and the li/la pseudo ladders (llvm-mc for li, a semantics check for
// la). Generators are arb/loong64 on top of arch/loong64 constructors.
// The common part of the suite is property_test.go; the larch suffix
// (not loong64!) mirrors the a64/rv convention of property_test.go: the
// _<GOARCH> suffix of the file name is a build constraint, and a real
// GOARCH name would keep the file off this darwin/arm64 host.

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	mrnd "math/rand/v2"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

	ohsnap "github.com/okneniz/oh-snap"
	parsecbytes "github.com/okneniz/parsec/bytes"
	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/arb"
	larch "github.com/okneniz/assembly/arb/loong64"
	"github.com/okneniz/assembly/arch/loong64"
	"github.com/okneniz/assembly/asm/loong64/pseudo"
	"github.com/okneniz/assembly/disasm"
	"github.com/okneniz/assembly/file"
	"github.com/okneniz/assembly/tests/cmd/objdump"
)

// laText - normalized ObjDump text of a loong64 instruction.
func laText(in loong64.Instr) string {
	return objdump.StripComments(objdump.Normalize(in.ObjDump(disasm.DefaultViewCtx())))
}

// laCanon - the decode-only aliases back into the base forms the
// assembler accepts. The rdcnt* names are objdump notation (llvm prints
// them too), but neither the syntax layer nor pseudo has a spelling for
// them; the base form assembles into the same word, so the round trip
// stays honest. Everything else passes through untouched. Applied per
// line (the input is a whole program in the list properties).
func laCanon(src string) string {
	lines := strings.Split(src, "\n")
	for i, ln := range lines {
		fields := strings.Fields(ln)
		switch {
		case len(fields) == 2 && fields[0] == "rdcntid.w":
			lines[i] = "rdtimel.w $zero, " + fields[1]
		case len(fields) == 2 && fields[0] == "rdcntvl.w":
			lines[i] = "rdtimel.w " + fields[1] + ", $zero"
		case len(fields) == 2 && fields[0] == "rdcntvh.w":
			lines[i] = "rdtimeh.w " + fields[1] + ", $zero"
		}
	}

	return strings.Join(lines, "\n")
}

// laAssemblesTo - bytes from assembling the text (false on assembly
// error); the decode-only aliases are canonicalized first.
func laAssemblesTo(t *testing.T, src string) ([]byte, bool) {
	t.Helper()
	res, errs := pseudo.Assemble(laCanon(src), propAddr)
	if len(errs) != 0 {
		t.Logf("%q: assemble: %v", src, errs)
		return nil, false
	}

	return res.Sections[0].Data, true
}

// laBytesOf - the instruction word (LA64 words are fixed 32-bit).
func laBytesOf(t *testing.T, in loong64.Instr) ([]byte, bool) {
	t.Helper()
	var buf bytes.Buffer
	if _, err := in.Encode(&buf, propAddr); err != nil {
		t.Logf("%s: Encode: %v", in.ObjDump(disasm.DefaultViewCtx()), err)
		return nil, false
	}

	return buf.Bytes(), true
}

// laEnc - the encoding context of the "bytes" property: a fixed address
// (PC-relative forms), no symbols. It equals the branch families' base
// (larch.BranchBase), so every generated target encodes.
type laEnc struct {
	addr uint64
}

// laEncodeAll - encodes a list sequentially starting at propAddr (each
// instruction at propAddr + 4*i; branch targets keep the family-span
// slack above the ±131068 reach of the narrowest forms).
func laEncodeAll(t *testing.T, ins []loong64.Instr) ([]byte, bool) {
	t.Helper()
	var buf bytes.Buffer
	addr := uint64(propAddr)
	for _, in := range ins {
		n, err := in.Encode(&buf, addr)
		if err != nil {
			t.Logf("encode: %v", err)
			return nil, false
		}

		addr += uint64(n)
	}

	return buf.Bytes(), true
}

// laBytesRoundTrip - the "bytes" property: the law enc∘dec∘enc == enc
// (RoundTrip) in the encoding context propAddr without symbols - the
// word is stable after the round trip.
func laBytesRoundTrip(t *testing.T, in loong64.Instr) bool {
	t.Helper()
	return RoundTrip[laEnc, loong64.Instr, []byte](
		laEnc{addr: propAddr},
		func(_ laEnc, x loong64.Instr) ([]byte, bool) {
			return laBytesOf(t, x)
		},
		func(ctx laEnc, b []byte) (loong64.Instr, bool) {
			back, err := loong64.Parse(ctx.addr)(parsecbytes.Buffer(b))
			if err != nil {
				t.Logf("decode: %v", err)
				return nil, false
			}

			if len(back) != 1 {
				t.Logf("% x: decode = %d instr", b, len(back))
				return nil, false
			}

			return back[0], true
		},
		bytes.Equal,
	)(in)
}

// laTextRoundTrip - the "text" property: the RoundTrip law in the
// DefaultViewCtx context - a fixed point of the "text → assembly →
// decode" cycle (the alias forms collapse into a single canon: move and
// its base form print identically). The input is canonicalized through
// bytes and decode: the fixed point is sought for the decoder's text.
func laTextRoundTrip(t *testing.T, in loong64.Instr) bool {
	t.Helper()
	b, ok := laBytesOf(t, in)
	if !ok {
		return false
	}

	d1, err := loong64.Parse(propAddr)(parsecbytes.Buffer(b))
	if err != nil {
		t.Logf("decode: %v", err)
		return false
	}

	if len(d1) != 1 {
		t.Logf("% x: decode = %d instr", b, len(d1))
		return false
	}

	return RoundTrip[disasm.ViewCtx, loong64.Instr, string](
		disasm.DefaultViewCtx(),
		func(_ disasm.ViewCtx, y loong64.Instr) (string, bool) {
			return laText(y), true
		},
		func(_ disasm.ViewCtx, src string) (loong64.Instr, bool) {
			data, ok := laAssemblesTo(t, src)
			if !ok {
				return nil, false
			}

			d2, err := loong64.Parse(propAddr)(parsecbytes.Buffer(data))
			if err != nil {
				t.Logf("decode: %v", err)
				return nil, false
			}

			if len(d2) != 1 {
				t.Logf("%q: decode = %d instr", src, len(d2))
				return nil, false
			}

			return d2[0], true
		},
		func(a, b string) bool { return a == b },
	)(d1[0])
}

// laInstrParam - parameters of any loong64 family.
type laInstrParam interface {
	Instr() loong64.Instr
	String() string
}

// laFamilyEntry - one loong64 family in the round-trip tables.
type laFamilyEntry struct {
	name string
	run  func(t *testing.T, rnd *mrnd.Rand)
}

// newLaFamily - a loong64 family entry: closes over the generic
// instantiation of the properties (families have different parameter
// types, Arbitrary is invariant - see newRvFamily in property_rv_test.go;
// each property is a separate named subtest, the check runs over the
// parameters for the sake of shrinking).
func newLaFamily[P laInstrParam](
	name string,
	mk func(rnd *mrnd.Rand) ohsnap.Arbitrary[P],
) laFamilyEntry {
	return laFamilyEntry{
		name: name,
		run: func(t *testing.T, rnd *mrnd.Rand) {
			t.Helper()
			t.Run("bytes", func(t *testing.T) {
				ohsnap.Check(t, 300, mk(rnd), func(p P) bool {
					return laBytesRoundTrip(t, p.Instr())
				})
			})

			t.Run("text", func(t *testing.T) {
				ohsnap.Check(t, 300, mk(rnd), func(p P) bool {
					return laTextRoundTrip(t, p.Instr())
				})
			})
		},
	}
}

// laFamilies - the round-trip family table (every generator family of
// arb/loong64; the operandless forms are differential base words).
func laFamilies() []laFamilyEntry {
	return []laFamilyEntry{
		newLaFamily("Alu3R", larch.Alu3R),
		newLaFamily("Alu2R", larch.Alu2R),
		newLaFamily("AluImm12", larch.AluImm12),
		newLaFamily("AluUImm12", larch.AluUImm12),
		newLaFamily("AluImm16", larch.AluImm16),
		newLaFamily("Imm20", larch.Imm20Instr),
		newLaFamily("Code15", larch.Code15Instr),
		newLaFamily("Branch2", larch.Branch2),
		newLaFamily("Branch1", larch.Branch1),
		newLaFamily("Jump", larch.Jump),
		newLaFamily("Jirl", larch.Jirl),
		newLaFamily("LdSt", larch.LdSt),
		newLaFamily("Ldptr", larch.Ldptr),
		newLaFamily("LdxStx", larch.LdxStx),
		newLaFamily("LdAcq", larch.LdAcq),
		newLaFamily("Hints", larch.Hints),
		newLaFamily("Preldx", larch.Preldx),
		newLaFamily("ShiftW", larch.ShiftW),
		newLaFamily("ShiftD", larch.ShiftD),
		newLaFamily("FieldW", larch.FieldW),
		newLaFamily("FieldD", larch.FieldD),
		newLaFamily("Alsl", larch.Alsl),
		newLaFamily("BytepickW", larch.BytepickW),
		newLaFamily("BytepickD", larch.BytepickD),
		newLaFamily("Atomics", larch.Atomics),
		newLaFamily("CsrRW", larch.CsrRW),
		newLaFamily("CsrXchg", larch.CsrXchg),
		newLaFamily("IoCsr", larch.IoCsr),
		newLaFamily("Lddir", larch.Lddir),
		newLaFamily("Ldpte", larch.Ldpte),
		newLaFamily("Invtlb", larch.Invtlb),
	}
}

// TestPropertyLoongSingleInstrRoundTrip - round trip of each loong64
// family: the word (bytes) and the canonical text are fixed points of
// encode∘decode; the pc-relative forms round-trip against their own pc
// (propAddr - the branch families generate targets around it).
func TestPropertyLoongSingleInstrRoundTrip(t *testing.T) {
	for _, f := range laFamilies() {
		t.Run(f.name, func(t *testing.T) {
			f.run(t, seedRnd(t))
		})
	}
}

// TestPropertyLoongAliasRoundTrip - the "alias" property (decode
// aliases): an instruction written via its alias (or the base form of a
// decode-only alias) decodes into the alias text, and the alias text
// re-assembles into the same bytes.
func TestPropertyLoongAliasRoundTrip(t *testing.T) {
	for _, tc := range []struct {
		name string
		src  string
		want string
	}{
		{"nop", "andi $zero, $zero, 0", "nop"},
		{"move", "or $t0, $t1, $zero", "move $t0, $t1"},
		{"move-pseudo", "move $t0, $t1", "move $t0, $t1"},
		{"ret", "ret", "ret"},
		{"ret-base", "jirl $zero, $ra, 0", "ret"},
		{"jr", "jr $t1", "jr $t1"},
		{"jr-base", "jirl $zero, $t1, 0", "jr $t1"},
		{"bltz", "bltz $t1, 8", "bltz $t1, 8"},
		{"bltz-base", "blt $t1, $zero, 8", "bltz $t1, 8"},
		{"bgez", "bgez $t1, 8", "bgez $t1, 8"},
		{"bgez-base", "bge $t1, $zero, 8", "bgez $t1, 8"},
		{"bgtz", "bgtz $t0, 8", "bgtz $t0, 8"},
		{"bgtz-base", "blt $zero, $t0, 8", "bgtz $t0, 8"},
		{"blez", "blez $t0, 8", "blez $t0, 8"},
		{"blez-base", "bge $zero, $t0, 8", "blez $t0, 8"},
		{"rdcntid", "rdtimel.w $zero, $t1", "rdcntid.w $t1"},
		{"rdcntvl", "rdtimel.w $t0, $zero", "rdcntvl.w $t0"},
		{"rdcntvh", "rdtimeh.w $t0, $zero", "rdcntvh.w $t0"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			res, errs := pseudo.Assemble(tc.src, propAddr)
			require.Empty(t, errs, "assemble %q", tc.src)
			data := res.Sections[0].Data

			back, err := loong64.Parse(propAddr)(parsecbytes.Buffer(data))
			require.NoError(t, err)
			require.Len(t, back, 1, "%q", tc.src)
			require.Equal(t, tc.want, laText(back[0]), "%q", tc.src)

			again, ok := laAssemblesTo(t, tc.want)
			require.True(t, ok, "re-assemble %q", tc.want)
			require.Equal(t, data, again, "alias text %q", tc.want)
		})
	}
}

// laInstrOf - a sampler function over a loong64 family generator (the
// generator is created once; see rvInstrOf in property_rv_test.go).
func laInstrOf[P laInstrParam](a ohsnap.Arbitrary[P]) func() loong64.Instr {
	return func() loong64.Instr {
		return a.Generate().Instr()
	}
}

// laPropFamilies - loong64 family generators for composition and the
// differential.
func laPropFamilies(rnd *mrnd.Rand) []func() loong64.Instr {
	return []func() loong64.Instr{
		laInstrOf(larch.Alu3R(rnd)),
		laInstrOf(larch.Alu2R(rnd)),
		laInstrOf(larch.AluImm12(rnd)),
		laInstrOf(larch.AluUImm12(rnd)),
		laInstrOf(larch.AluImm16(rnd)),
		laInstrOf(larch.Imm20Instr(rnd)),
		laInstrOf(larch.Code15Instr(rnd)),
		laInstrOf(larch.Branch2(rnd)),
		laInstrOf(larch.Branch1(rnd)),
		laInstrOf(larch.Jump(rnd)),
		laInstrOf(larch.Jirl(rnd)),
		laInstrOf(larch.LdSt(rnd)),
		laInstrOf(larch.Ldptr(rnd)),
		laInstrOf(larch.LdxStx(rnd)),
		laInstrOf(larch.LdAcq(rnd)),
		laInstrOf(larch.Hints(rnd)),
		laInstrOf(larch.Preldx(rnd)),
		laInstrOf(larch.ShiftW(rnd)),
		laInstrOf(larch.ShiftD(rnd)),
		laInstrOf(larch.FieldW(rnd)),
		laInstrOf(larch.FieldD(rnd)),
		laInstrOf(larch.Alsl(rnd)),
		laInstrOf(larch.BytepickW(rnd)),
		laInstrOf(larch.BytepickD(rnd)),
		laInstrOf(larch.Atomics(rnd)),
		laInstrOf(larch.CsrRW(rnd)),
		laInstrOf(larch.CsrXchg(rnd)),
		laInstrOf(larch.IoCsr(rnd)),
		laInstrOf(larch.Lddir(rnd)),
		laInstrOf(larch.Ldpte(rnd)),
		laInstrOf(larch.Invtlb(rnd)),
	}
}

// TestPropertyLoongBytesRoundTripList - the "bytes" property for a hand
// list of llvm-verified programs: the source (pseudo li/la/call/tail and
// real instructions with labels) assembles into exactly these words.
func TestPropertyLoongBytesRoundTripList(t *testing.T) {
	for _, tc := range []struct {
		name  string
		base  uint64
		src   string
		words []uint32
	}{
		{
			"alu", 0, "add.w $t0, $t1, $t2",
			[]uint32{0x001039ac},
		},
		{
			"imm", 0, "addi.w $t0, $t1, -16\nandi $t0, $t1, 0xf0f",
			[]uint32{0x02bfc1ac, 0x037c3dac},
		},
		{
			"lu-st", 0, "lu12i.w $t0, 0x12345\nld.d $t0, $t1, 8\nst.w $t0, $t1, -8",
			[]uint32{0x142468ac, 0x28c021ac, 0x29bfe1ac},
		},
		{
			"shift-field", 0, "slli.d $t0, $t1, 3\nbstrins.w $t0, $t1, 5, 3\nalsl.w $t0, $t1, $t2, 3",
			[]uint32{0x00410dac, 0x00650dac, 0x000539ac},
		},
		{
			"priv", 0, "csrrd $t0, 5\ntlbsrch\ndbar 0",
			[]uint32{0x0400140c, 0x06482800, 0x38720000},
		},
		{
			"pseudo-mix", 0, "nop\nmove $t0, $t1\nadd.w $t0, $t0, $t2",
			[]uint32{0x03400000, 0x001501ac, 0x0010398c},
		},
		{
			"bltz", 0, "bltz $t1, 8",
			[]uint32{0x600009a0},
		},
		{
			"bgez", 0, "bgez $t1, 8",
			[]uint32{0x640009a0},
		},
		{
			"bgtz", 0, "bgtz $t0, 8",
			[]uint32{0x6000080c},
		},
		{
			"blez", 0, "blez $t0, 8",
			[]uint32{0x6400080c},
		},
		{
			"labels", 0x90000000, "beqz $t1, 1f\n1: b 1b",
			[]uint32{0x400005a0, 0x50000000},
		},
		{
			"labels-fwd-back", 0, "1: b 1f\nb 1b\n1: b 1b",
			[]uint32{0x50000800, 0x53ffffff, 0x50000000},
		},
		{
			"call", 0, "call 8",
			[]uint32{0x54000800},
		},
		{
			"tail", 0, "tail 8",
			[]uint32{0x50000800},
		},
		{
			"call-label", 0x90000000, "call 1f\n1: nop",
			[]uint32{0x54000400, 0x03400000},
		},
		{
			"li-short", 0, "li.w $t0, -4096",
			[]uint32{0x15ffffec},
		},
		{
			"li-ori", 0, "li.d $t0, 4095",
			[]uint32{0x03bffc0c},
		},
		{
			"li-pair", 0, "li.d $t0, 0x12345678",
			[]uint32{0x142468ac, 0x0399e18c},
		},
		{
			"li-neg-pair", 0, "li.d $t0, -414771",
			[]uint32{0x15fff34c, 0x03af358c},
		},
		{
			"li-triple", 0, "li.d $t0, 0x500000000",
			[]uint32{0x1400000c, 0x160000ac, 0x0300018c},
		},
		{
			"li-worst", 0, "li.d $t0, 0x123456789abcdef0",
			[]uint32{0x153579ac, 0x03bbc18c, 0x168acf0c, 0x03048d8c},
		},
		{
			"la", 0x90000000, "la $t0, .+0x1024",
			[]uint32{0x1a00002c, 0x02c0918c},
		},
		{
			"la-neg-lo", 0x90000000, "la $t0, .+0xf80",
			[]uint32{0x1a00002c, 0x02fe018c},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			res, errs := pseudo.Assemble(tc.src, tc.base)
			require.Empty(t, errs, "assemble %q", tc.src)
			data := res.Sections[0].Data
			require.Len(t, data, 4*len(tc.words), tc.name)

			for i, want := range tc.words {
				require.Equal(
					t,
					want,
					binary.LittleEndian.Uint32(data[i*4:]),
					"%s word %d",
					tc.name,
					i,
				)
			}
		})
	}
}

// TestPropertyLoongTextRoundTripList - the "text" property for a hand
// list of programs: the source assembles, the bytes decode line by line
// into ObjDump texts, and the joined texts (the same $-names; the
// decode-only aliases canonicalized) re-assemble into the same bytes.
func TestPropertyLoongTextRoundTripList(t *testing.T) {
	for _, tc := range []struct {
		name string
		src  string
	}{
		{"alu", "add.w $t0, $t1, $t2\naddi.d $a0, $a1, -1\nsltui $s0, $s1, 2047"},
		{"alias", "nop\nmove $t0, $t1\nnot $t0, $t1\nret"},
		{"rdcnt", "rdtimel.w $zero, $t1\nrdtimel.w $t0, $zero\nrdtimeh.w $t0, $zero"},
		{"branch-labels", "1: beq $t0, $t1, 1b\nb 1b"},
		{"branch-alias-labels", "1: bltz $t1, 1f\nbgez $t2, 1b\n1: blez $t3, 1b"},
		{"jump-reg", "jirl $ra, $t0, 16\njirl $t0, $t1, -20"},
		{"li-mix", "li.d $t0, 0x12345678\nadd.w $t0, $t0, $t2\nli.w $t1, -4096"},
		{"la-mix", "la $t0, .+0x1024\nld.d $t1, $t0, 0"},
		{"mem", "ldx.bu $a0, $a1, $a2\nstle.d $t0, $t1, $t2\npreld 3, $sp, -8"},
		{"shift-bits", "slli.d $t0, $t1, 3\nbstrpick.d $t0, $t1, 33, 17\nbytepick.w $t0, $t1, $t2, 2"},
		{"atomics", "amadd.w $t0, $t1, $t2\namswap_db.d $a0, $a1, $a2\nsc.q $t0, $t1, $t2"},
		{"priv", "csrrd $t0, 5\ncsrwr $t1, 0x123\ncsrxchg $t2, $t3, 16383\ninvtlb 0, $t0, $t1"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			res, errs := pseudo.Assemble(tc.src, propAddr)
			require.Empty(t, errs, "assemble %q", tc.src)
			data := res.Sections[0].Data

			back, err := loong64.Parse(propAddr)(parsecbytes.Buffer(data))
			require.NoError(t, err)
			require.NotEmpty(t, back, tc.name)

			texts := make([]string, len(back))
			for i := range back {
				texts[i] = laText(back[i])
			}

			again, ok := laAssemblesTo(t, strings.Join(texts, "\n"))
			require.True(t, ok, "%s: re-assemble %q", tc.name, texts)
			require.Equal(t, data, again, "%s: bytes of %q", tc.name, texts)

			back2, err := loong64.Parse(propAddr)(parsecbytes.Buffer(again))
			require.NoError(t, err)
			require.Len(t, back2, len(texts), tc.name)
			for i := range back2 {
				require.Equal(t, texts[i], laText(back2[i]), "%s[%d]", tc.name, i)
			}
		})
	}
}

// TestPropertyLoongSeqRoundTripList - the "bytes" property for random
// compositions: a generated list encodes, decodes line by line, and
// encodes again into the same bytes (variable contents, fixed 4-byte
// width).
func TestPropertyLoongSeqRoundTripList(t *testing.T) {
	rnd := seedRnd(t)
	seq := arb.Seq(rnd, laPropFamilies(rnd))
	ohsnap.Check(t, 100, seq, func(ins []loong64.Instr) bool {
		raw, ok := laEncodeAll(t, ins)
		if !ok {
			return false
		}

		back, err := loong64.Parse(propAddr)(parsecbytes.Buffer(raw))
		if err != nil {
			t.Logf("decode: %v", err)
			return false
		}

		if len(back) != len(ins) {
			t.Logf("decoded %d of %d", len(back), len(ins))
			return false
		}

		raw2, ok := laEncodeAll(t, back)
		if !ok {
			return false
		}

		if !bytes.Equal(raw2, raw) {
			t.Logf("re-encode: % x ≠ % x", raw2, raw)
			return false
		}

		return true
	})
}

// TestPropertyLoongDecodeRobustness - arbitrary words do not crash the
// decoder; truncated buffers decode into a prefix.
func TestPropertyLoongDecodeRobustness(t *testing.T) {
	rnd := seedRnd(t)
	t.Run("words", func(t *testing.T) {
		ohsnap.Check(t, 100000, arb.Word(rnd), func(w uint32) bool {
			ok := true
			func() {
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("panic at %#08x: %v", w, r)
						ok = false
					}
				}()
				data := binary.LittleEndian.AppendUint32(nil, w)
				ins, err := loong64.Parse(propAddr)(parsecbytes.Buffer(data))
				if err != nil {
					t.Errorf("parse %#08x: %v", w, err)
					ok = false
					return
				}

				// a whole word always decodes (Unknown as the fallback)
				if len(ins) != 1 {
					t.Errorf("word %#08x: %d instructions", w, len(ins))
					ok = false
					return
				}

				for _, in := range ins {
					_ = in.ObjDump(disasm.DefaultViewCtx())

					if _, err := in.MarshalJSON(); err != nil {
						t.Logf("word %#08x: MarshalJSON: %v", w, err)
						ok = false
					}

					if n := in.Len(); n != 4 {
						t.Logf("word %#08x: Len = %d", w, n)
						ok = false
					}
				}
			}()
			return ok
		})
	})

	t.Run("truncated", func(t *testing.T) {
		for _, w := range []uint32{0x001039ac, 0x03400000, 0xffffffff, 0x4c000020} {
			data := binary.LittleEndian.AppendUint32(nil, w)
			for n := 1; n < 4; n++ {
				ins, err := loong64.Parse(propAddr)(parsecbytes.Buffer(data[:n]))
				require.NoError(t, err, "%#08x[:%d]", w, n)
				require.Empty(t, ins, "%#08x[:%d]: a partial word is not an instruction", w, n)
			}

			for n := 5; n < 8; n++ {
				ins, err := loong64.Parse(propAddr)(parsecbytes.Buffer(data[:n]))
				require.NoError(t, err, "%#08x[:%d]", w, n)
				require.Len(t, ins, 1, "%#08x[:%d]: the whole-word prefix", w, n)
			}
		}
	})
}

// laBranchMnemonics - the forms whose LAST operand is a pc-relative
// branch target: both llvm-objdump and GNU objdump print it as the raw
// byte OFFSET (with the absolute address as a comment/annotation), our
// frozen decoder prints the absolute TARGET - a known arch deviation
// (reported). The differential compares these lines field by field and
// checks the target arithmetic instead (ours == addr + objdump's).
var laBranchMnemonics = map[string]bool{
	"beq": true, "bne": true, "blt": true, "bge": true,
	"bltu": true, "bgeu": true,
	"beqz": true, "bnez": true,
	"b": true, "bl": true,
	"bltz": true, "bgez": true, "bgtz": true, "blez": true,
}

// laIsHexRunes - every rune of s is a hex digit.
func laIsHexRunes(s string) bool {
	for _, r := range s {
		if !strings.ContainsRune("0123456789abcdef", r) {
			return false
		}
	}

	return s != ""
}

// laInstrTail - the instruction fields of a normalized objdump line
// (the mnemonic and operands, without the address and code columns).
// llvm prints the LoongArch code column as four 2-digit bytes, GNU as a
// single 8-digit word.
func laInstrTail(line string) []string {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return nil
	}

	if len(fields[1]) == 8 && laIsHexRunes(fields[1]) {
		return fields[2:]
	}

	if len(fields) >= 6 &&
		len(fields[1]) == 2 && len(fields[2]) == 2 &&
		len(fields[3]) == 2 && len(fields[4]) == 2 &&
		laIsHexRunes(fields[1]) && laIsHexRunes(fields[2]) &&
		laIsHexRunes(fields[3]) && laIsHexRunes(fields[4]) {
		return fields[5:]
	}

	return nil
}

// laLineEqual - the differential comparison of one address: our
// instruction text against the instruction tail of the objdump line.
// Strict equality, except the pc-relative branch forms
// (laBranchMnemonics): the fields before the target must be equal and
// the targets must agree arithmetically (our absolute == addr +
// objdump's offset).
func laLineEqual(addr uint64, ourInstr, objLine string) bool {
	tail := laInstrTail(objLine)
	if tail == nil {
		return false
	}

	if ourInstr == strings.Join(tail, " ") {
		return true
	}

	of := strings.Fields(ourInstr)
	if len(of) == 0 || len(of) != len(tail) || !laBranchMnemonics[of[0]] {
		return false
	}

	for i := range len(of) - 1 {
		if of[i] != tail[i] {
			return false
		}
	}

	ourAbs, err1 := strconv.ParseInt(of[len(of)-1], 10, 64)
	objOff, err2 := strconv.ParseInt(tail[len(tail)-1], 10, 64)
	if err1 != nil || err2 != nil {
		return false
	}

	return ourAbs == int64(addr)+objOff
}

// TestPropertyLoongVsObjdump - the differential: words produced by the
// loong64 generators and the base word of every decode-table entry are
// disassembled identically by us and by objdump (llvm: Apple's objdump
// has no LoongArch target, objdump.Run falls through to homebrew). Two
// known deviations are normalized, not silently ignored: the branch
// target print (relative vs absolute, laLineEqual) is counted in its own
// bucket and reported; llvm prints the LoongArch code column as bytes
// (unlike riscv, where it prints a hex word), so our line is rendered in
// the byte style.
func TestPropertyLoongVsObjdump(t *testing.T) {
	const perFamily = 48
	const threshold = 90.0

	rnd := seedRnd(t)

	// every generated instruction is encoded at one and the same pc
	// (propAddr): the targets of the pc-relative families are generated
	// around it, and the comparison is word+address based anyway - both
	// objdump and our decoder recompute a branch target from the address
	// the word sits at.
	var buf bytes.Buffer
	for _, gen := range laPropFamilies(rnd) {
		for range perFamily {
			_, err := gen().Encode(&buf, propAddr)
			require.NoError(t, err)
		}
	}

	// the base word of every decode-table mnemonic: the empty forms
	// (tlbsrch and friends) ride along here - they have no generators.
	for _, m := range loong64.Mnemonics() {
		buf.Write(binary.LittleEndian.AppendUint32(nil, loong64.EncodingWord(m)))
	}

	code := buf.Bytes()
	path := writeLoongELF(t, code)
	out, err := objdump.Run(context.Background(), objdump.Args("ELF", 0, path))
	if err != nil {
		t.Skipf("no objdump for ELF/loongarch64: %v", err)
	}

	objLines := objdump.ParseByAddr(string(out))
	require.NotEmpty(t, objLines)

	// the comparison is per instruction text (the objdump line tail):
	// llvm prints the LoongArch code column as bytes and GNU as a word,
	// and the words are ours to begin with - both decoders read them
	// from the ELF we just wrote
	matched, branchRel, mismatched, notInOurs := 0, 0, 0, 0
	var samples []string
	for addr, objLine := range objLines {
		off := int(addr) - propAddr
		if off < 0 || off+4 > len(code) {
			notInOurs++
			continue
		}

		ins, err := loong64.Parse(addr)(parsecbytes.Buffer(code[off:]))
		if err != nil || len(ins) == 0 {
			notInOurs++
			continue
		}

		ourText := laText(ins[0])
		objText := objdump.StripComments(objLine)
		if laLineEqual(addr, ourText, objText) {
			tail := laInstrTail(objText)
			if tail != nil && ourText == strings.Join(tail, " ") {
				matched++
			} else {
				branchRel++
			}
		} else {
			mismatched++
			if len(samples) < 10 {
				samples = append(samples, fmt.Sprintf(
					"0x%x (%08x)\n    ours:    %s\n    objdump: %s",
					addr,
					binary.LittleEndian.Uint32(code[off:off+4]),
					ourText,
					objText,
				))
			}
		}
	}

	total := matched + branchRel + mismatched
	pct := 0.0
	if total > 0 {
		pct = float64(matched+branchRel) * 100 / float64(total)
	}

	t.Logf(
		"vs objdump: %d matched (%.2f%%), %d branch prints normalized"+
			" (known relative-vs-absolute deviation, reported),"+
			" %d mismatched, %d foreign addresses",
		matched, pct, branchRel, mismatched, notInOurs)
	for _, sm := range samples {
		t.Log(sm)
	}

	require.GreaterOrEqual(t, pct, threshold)
}

// writeLoongELF - minimal ELF64 LE loongarch (e_machine = EM_LOONGARCH,
// e_flags = double-float ABI, as real LoongArch toolchains emit):
// .text @ 0x1000 (the shared writeObjELF writer, machine-parameterized).
func writeLoongELF(t *testing.T, code []byte) string {
	t.Helper()
	const emLoongArch = 258
	const efLoongArchDoubleFloat = 0x2
	path, err := writeObjELF(t, code, emLoongArch, efLoongArchDoubleFloat)
	require.NoError(t, err)
	return path
}

// llvmMcPath - the llvm-mc binary (the assembler oracle), probed the
// objdump.Run way: fixed candidates first, then PATH. Empty when absent
// (the li differential skips).
func llvmMcPath() string {
	candidates := []string{
		"/opt/homebrew/opt/llvm/bin/llvm-mc",
		"/opt/homebrew/bin/llvm-mc",
		"/usr/local/opt/llvm/bin/llvm-mc",
		"/usr/local/bin/llvm-mc",
		"llvm-mc",
	}
	for _, c := range candidates {
		if _, err := exec.LookPath(c); err == nil {
			return c
		}
	}

	return ""
}

// mcAssemble - bytes of the .text llvm-mc assembles from src
// (loongarch64 object); the test fails on an llvm-mc error.
func mcAssemble(t *testing.T, mc, src string) []byte {
	t.Helper()
	path := t.TempDir() + "/mc.o"
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, mc, "-filetype=obj", "--triple=loongarch64", "-o", path)
	cmd.Stdin = strings.NewReader(src)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("llvm-mc %q: %v: %s", src, err, out)
		return nil
	}

	f, err := file.Detect(path)
	if err != nil {
		t.Fatalf("detect %s: %v", path, err)
		return nil
	}

	sec, err := f.CodeSection()
	if err != nil {
		t.Fatalf(".text of %s: %v", path, err)
		return nil
	}

	return sec.Data
}

// knownLiDeviations - the llvm-mc ladder deviations of the frozen
// asm/loong64/pseudo li expansion, pinned per mnemonic+value so that any
// NEW deviation fails the differential while the known ones stay
// visible. Both encodings load the same constant; the classes:
//   - "addi-rung": 0..2047 - ours takes the addi.w rung first, llvm-mc
//     takes ori (one word either way);
//   - "zero-group": bits above 32 with an empty low half - ours always
//     emits the fixed 3-word chain (the lu12i.w 0 of it is redundant),
//     llvm-mc drops the empty groups;
//   - "sign-trick": the int64 edges - llvm-mc exploits the $zero base
//     and sign-extension (a single lu52i.d for the minimum, addi.w -1 +
//     lu52i.d for the maximum), ours the fixed 3/4-word chains.
var knownLiDeviations = map[string]map[int64]string{
	// The remaining llvm deviations are compiler-grade chain tricks:
	// llvm drops a redundant leading zero word (ori 0) and exploits an
	// all-ones addi.w base for the int64 sign edges.
	"li.d": {
		0x500000000:         "zero-group",
		0x7fffffffffffffff:  "sign-trick",
		-0x8000000000000000: "sign-trick",
	},
}

// TestPropertyLoongLiVsLlvm - the li.d/li.w ladders against llvm-mc:
// every ladder rung (and the negatives, and the 64-bit constants)
// assembles into exactly the bytes llvm-mc chooses, except the pinned
// deviations above (logged, reported - the pseudo layer is frozen).
// Skipped when no llvm-mc is installed.
func TestPropertyLoongLiVsLlvm(t *testing.T) {
	mc := llvmMcPath()
	if mc == "" {
		t.Skip("no llvm-mc found on this host")
	}

	for _, tc := range []struct {
		mnem   string
		values []int64
	}{
		{"li.d", []int64{
			-2048, -16, -1, 0, 1, 2047, // addi.w rung (llvm: ori for 0..2047)
			4095, 4096, // ori rung edge
			0x12345000, 0x12345678, // lu12i.w rungs
			-4096, -414771, // negative pages
			0x500000000,                             // the 3-word rung (bits above 32)
			0x123456789abcdef0,                      // the worst 4-word chain
			0x7fffffffffffffff, -0x8000000000000000, // the int64 edges
		}},
		{"li.w", []int64{
			-2048, -16, -1, 0, 1, 2047,
			4095, 4096, 0x12345000, 0x12345678, -4096, -0x80000000,
		}},
	} {
		t.Run(tc.mnem, func(t *testing.T) {
			for _, v := range tc.values {
				src := fmt.Sprintf("%s $t0, %#x", tc.mnem, v)
				res, errs := pseudo.Assemble(src, propAddr)
				require.Empty(t, errs, "assemble %q", src)

				want := mcAssemble(t, mc, src)
				if class, known := knownLiDeviations[tc.mnem][v]; known &&
					!bytes.Equal(want, res.Sections[0].Data) {
					t.Logf(
						"KNOWN llvm deviation (%s, reported - pseudo frozen): %s of %#x"+
							" - ours % x, llvm % x",
						class, tc.mnem, v, res.Sections[0].Data, want,
					)
					continue
				}

				require.Equal(
					t,
					want,
					res.Sections[0].Data,
					"%s of %#x: our ladder ≠ llvm-mc ladder",
					tc.mnem,
					v,
				)
			}
		})
	}
}

// laSext - v (already masked to its field) as a signed n-bit value (the
// la pair field extraction below).
func laSext(v uint32, bits int) int64 {
	return int64(int32(v<<(32-bits))) >> (32 - bits)
}

// TestPropertyLoongLaSemantics - the la pair as a self-consistency
// property: our own two words decode back into the pcalau12i+addi.d
// texts, and the pair arithmetic (page + si20<<12 + si12) reconstructs
// the target exactly (llvm requires a bare symbol for la, so the pair is
// checked against its own decoding instead of an llvm-mc diff; the words
// themselves are llvm-verified in the BytesRoundTripList hand table
// above).
func TestPropertyLoongLaSemantics(t *testing.T) {
	for _, tc := range []struct {
		name   string
		offset int64
	}{
		{"zero", 0},
		{"in-page", 0x123},
		{"page", 0x1000},
		{"page-cross", 0x1024},
		{"neg-lo", 0xf80}, // the sext12 carry goes into hi
		{"neg-page", -0x4000},
		{"far", 0x7fff000},
	} {
		t.Run(tc.name, func(t *testing.T) {
			src := fmt.Sprintf("la $t0, .+%#x", tc.offset)
			if tc.offset < 0 {
				src = fmt.Sprintf("la $t0, .-%#x", -tc.offset)
			}

			res, errs := pseudo.Assemble(src, propAddr)
			require.Empty(t, errs, "assemble %q", src)
			data := res.Sections[0].Data
			require.Len(t, data, 8, tc.name)

			back, err := loong64.Parse(propAddr)(parsecbytes.Buffer(data))
			require.NoError(t, err)
			require.Len(t, back, 2, tc.name)
			require.True(t, strings.HasPrefix(laText(back[0]), "pcalau12i $t0,"), tc.name)
			require.True(t, strings.HasPrefix(laText(back[1]), "addi.d $t0, $t0,"), tc.name)

			// the pair arithmetic: page + (si20 << 12) + si12 == target
			hi := laSext(binary.LittleEndian.Uint32(data[0:4])>>5&0xfffff, 20)
			lo := laSext(binary.LittleEndian.Uint32(data[4:8])>>10&0xfff, 12)
			want := int64(propAddr) + tc.offset
			got := int64(propAddr&^0xfff) + hi<<12 + lo
			require.Equal(t, want, got, "%s: the pair reconstructs the target", tc.name)
		})
	}

	// the out-of-range target refuses to assemble (the pair spans ±2 GiB)
	_, errs := pseudo.Assemble("la $t0, .+0x100000000", propAddr)
	require.NotEmpty(t, errs, "la beyond the pcalau12i+addi.d reach")
}
