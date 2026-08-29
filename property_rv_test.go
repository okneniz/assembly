package assembly_test

// Property tests (oh-snap), riscv: round trip of single instructions and
// compositions (variable length 2/4, canonicalization of pseudo-forms),
// decoder robustness, differential against the real objdump. Generators are
// arb/riscv on top of arch/riscv constructors. The common part of the suite
// is property_test.go.

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	mrnd "math/rand/v2"
	"strings"
	"testing"

	ohsnap "github.com/okneniz/oh-snap"
	parsecbytes "github.com/okneniz/parsec/bytes"
	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/arb"
	rv "github.com/okneniz/assembly/arb/riscv"
	"github.com/okneniz/assembly/arch/riscv"
	"github.com/okneniz/assembly/asm/riscv/pseudo"
	"github.com/okneniz/assembly/disasm"
	"github.com/okneniz/assembly/tests/cmd/objdump"
	"github.com/okneniz/assembly/text"
)

// rvText - normalized ObjDump text of a riscv instruction.
func rvText(in riscv.Instr) string {
	return objdump.StripComments(objdump.Normalize(in.ObjDump(disasm.DefaultViewCtx())))
}

// rvAssemblesTo - bytes from assembling the text (false on assembly error).
func rvAssemblesTo(t *testing.T, src string) ([]byte, bool) {
	t.Helper()
	res, errs := pseudo.Assemble(src, propAddr)
	if len(errs) != 0 {
		t.Logf("%q: assemble: %v", src, errs)
		return nil, false
	}

	return res.Sections[0].Data, true
}

// rvBytesOf - instruction bytes (2 or 4 - RVC compression).
func rvBytesOf(t *testing.T, in riscv.Instr) ([]byte, bool) {
	t.Helper()
	var buf bytes.Buffer
	if _, err := in.Encode(&buf, propAddr, riscv.EncOpts{}); err != nil {
		t.Logf("%s: Encode: %v", in.ObjDump(disasm.DefaultViewCtx()), err)
		return nil, false
	}

	return buf.Bytes(), true
}

// rvEnc - the encoding context of the "bytes" property: a fixed address
// (PC-relative forms), unrestricted modes, no symbols.
type rvEnc struct {
	addr uint64
}

// rvEncodeAll - encodes a list sequentially starting at propAddr (the address
// of each is propAddr + bytes written; compression allowed, no symbols).
func rvEncodeAll(t *testing.T, ins []riscv.Instr) ([]byte, bool) {
	t.Helper()
	var buf bytes.Buffer
	addr := propAddr
	for _, in := range ins {
		n, err := in.Encode(&buf, uint64(addr), riscv.EncOpts{})
		if err != nil {
			t.Logf("encode: %v", err)
			return nil, false
		}

		addr += int(n)
	}

	return buf.Bytes(), true
}

// rvBytesRoundTrip - the "bytes" property: the law enc∘dec∘enc == enc
// (RoundTrip) in the encoding context propAddr without symbols - bytes are
// stable after the round trip.
func rvBytesRoundTrip(t *testing.T, in riscv.Instr) bool {
	t.Helper()
	return RoundTrip[rvEnc, riscv.Instr, []byte](
		rvEnc{addr: propAddr},
		func(_ rvEnc, x riscv.Instr) ([]byte, bool) {
			return rvBytesOf(t, x)
		},
		func(ctx rvEnc, b []byte) (riscv.Instr, bool) {
			back, err := riscv.Parse(ctx.addr)(parsecbytes.Buffer(b))
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

// rvTextRoundTrip - the "text" property: the RoundTrip law in the
// DefaultViewCtx context - a fixed point of the "text → assembly → decode"
// cycle (pseudo-forms collapse into a single canon: mv and its base form
// print identically). The input is canonicalized through bytes and decode:
// the fixed point is sought for the decoder's text, while structures from
// parsing may be pseudo-forms whose text is not fixed.
func rvTextRoundTrip(t *testing.T, in riscv.Instr) bool {
	t.Helper()
	b, ok := rvBytesOf(t, in)
	if !ok {
		return false
	}

	d1, err := riscv.Parse(propAddr)(parsecbytes.Buffer(b))
	if err != nil {
		t.Logf("decode: %v", err)
		return false
	}

	if len(d1) != 1 {
		t.Logf("% x: decode = %d instr", b, len(d1))
		return false
	}

	return RoundTrip[disasm.ViewCtx, riscv.Instr, string](
		disasm.DefaultViewCtx(),
		func(_ disasm.ViewCtx, y riscv.Instr) (string, bool) {
			return rvText(y), true
		},
		func(_ disasm.ViewCtx, src string) (riscv.Instr, bool) {
			data, ok := rvAssemblesTo(t, src)
			if !ok {
				return nil, false
			}

			d2, err := riscv.Parse(propAddr)(parsecbytes.Buffer(data))
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

// rvInstrParam - parameters of any riscv family.
type rvInstrParam interface {
	Instr() riscv.Instr
}

// rvFamilyEntry - one riscv family in the TestPropertyRiscvSingleInstrRoundTrip table.
type rvFamilyEntry struct {
	name string
	run  func(t *testing.T, rnd *mrnd.Rand)
}

// newRvFamily - a riscv family entry: closes over the generic instantiation
// of the properties (families have different parameter types, Arbitrary is
// invariant - see newPropFamily in property_a64_test.go; each property is a
// separate named subtest, the check runs over the parameters for the sake of
// shrinking).
func newRvFamily[P rvInstrParam](
	name string,
	mk func(rnd *mrnd.Rand) ohsnap.Arbitrary[P],
) rvFamilyEntry {
	return rvFamilyEntry{
		name: name,
		run: func(t *testing.T, rnd *mrnd.Rand) {
			t.Helper()
			t.Run("bytes", func(t *testing.T) {
				ohsnap.Check(t, 300, mk(rnd), func(p P) bool {
					return rvBytesRoundTrip(t, p.Instr())
				})
			})

			t.Run("text", func(t *testing.T) {
				ohsnap.Check(t, 300, mk(rnd), func(p P) bool {
					return rvTextRoundTrip(t, p.Instr())
				})
			})
		},
	}
}

// TestPropertyRiscvSingleInstrRoundTrip - round trip of each riscv family.
func TestPropertyRiscvSingleInstrRoundTrip(t *testing.T) {
	families := []rvFamilyEntry{
		newRvFamily("Add", rv.Add),
		newRvFamily("Sub", rv.Sub),
		newRvFamily("Addi", rv.Addi),
		newRvFamily("Lui", rv.Lui),
		newRvFamily("Lw", rv.Lw),
		newRvFamily("Ld", rv.Ld),
		newRvFamily("Sw", rv.Sw),
		newRvFamily("Sd", rv.Sd),
	}
	for _, f := range families {
		t.Run(f.name, func(t *testing.T) {
			f.run(t, seedRnd(t))
		})
	}
}

// newRvAliasFamily - like newRvFamily, but the parameters are pinned
// (the alias condition of a pseudo-form); combinations invalid after pinning
// are outside the property. See newAliasFamily in property_a64_test.go.
func newRvAliasFamily[P rvInstrParam](
	name string,
	mk func(rnd *mrnd.Rand) ohsnap.Arbitrary[P],
	pin func(P) P,
) rvFamilyEntry {
	return rvFamilyEntry{
		name: name,
		run: func(t *testing.T, rnd *mrnd.Rand) {
			t.Helper()
			pinned := func(p P) (riscv.Instr, bool) {
				q := pin(p)
				if q.Instr() == nil {
					return nil, false
				}

				return q.Instr(), true
			}

			t.Run("bytes", func(t *testing.T) {
				ohsnap.Check(t, 300, mk(rnd), func(p P) bool {
					if in, ok := pinned(p); ok {
						return rvBytesRoundTrip(t, in)
					}

					return true
				})
			})

			t.Run("text", func(t *testing.T) {
				ohsnap.Check(t, 300, mk(rnd), func(p P) bool {
					if in, ok := pinned(p); ok {
						return rvTextRoundTrip(t, in)
					}

					return true
				})
			})
		},
	}
}

// TestPropertyRiscvAliasRoundTrip - the "alias" property (pseudo-forms):
// an instruction written as a pseudo-form survives the cycle
// struct → alias → struct without loss; each entry is its own subtest.
// nop is pinned entirely (the only instruction of the form) - the check
// degenerates into a constant one, but keeps the common framework in place.
func TestPropertyRiscvAliasRoundTrip(t *testing.T) {
	families := []rvFamilyEntry{
		newRvAliasFamily("mv", rv.Addi, func(p rv.AddiParams) rv.AddiParams {
			p.Imm = riscv.Imm12{}
			return p
		}),
		newRvAliasFamily("nop", rv.Addi, func(p rv.AddiParams) rv.AddiParams {
			p.Rd = riscv.Zero
			p.Rs1 = riscv.Zero
			p.Imm = riscv.Imm12{}
			return p
		}),
	}
	for _, f := range families {
		t.Run(f.name, func(t *testing.T) {
			f.run(t, seedRnd(t))
		})
	}
}

// rvInstrOf - a sampler function over a riscv family generator (the generator
// is created once; see instrOf in property_a64_test.go).
func rvInstrOf[P rvInstrParam](a ohsnap.Arbitrary[P]) func() riscv.Instr {
	return func() riscv.Instr {
		return ohsnap.First(a.Generate()).Instr()
	}
}

// rvPropFamilies - riscv family generators for composition and the differential.
func rvPropFamilies(rnd *mrnd.Rand) []func() riscv.Instr {
	return []func() riscv.Instr{
		rvInstrOf(rv.Add(rnd)),
		rvInstrOf(rv.Sub(rnd)),
		rvInstrOf(rv.Addi(rnd)),
		rvInstrOf(rv.Lui(rnd)),
		rvInstrOf(rv.Lw(rnd)),
		rvInstrOf(rv.Ld(rnd)),
		rvInstrOf(rv.Sw(rnd)),
		rvInstrOf(rv.Sd(rnd)),
	}
}

// TestPropertyRiscvBytesRoundTripList - the "bytes" property for a list of
// instructions: the list (variable length 2/4) is encoded, decoded line by
// line, and encoded again into the same bytes.
func TestPropertyRiscvBytesRoundTripList(t *testing.T) {
	rnd := seedRnd(t)
	seq := arb.Seq(rnd, rvPropFamilies(rnd))
	ohsnap.Check(t, 100, seq, func(ins []riscv.Instr) bool {
		raw, ok := rvEncodeAll(t, ins)
		if !ok {
			return false
		}

		buf := *bytes.NewBuffer(raw)

		back, err := riscv.Parse(propAddr)(parsecbytes.Buffer(buf.Bytes()))
		if err != nil {
			t.Logf("decode: %v", err)
			return false
		}

		if len(back) != len(ins) {
			t.Logf("decoded %d of %d", len(back), len(ins))
			return false
		}

		raw2, ok := rvEncodeAll(t, back)
		if !ok {
			return false
		}

		buf2 := *bytes.NewBuffer(raw2)

		if !bytes.Equal(buf2.Bytes(), buf.Bytes()) {
			t.Logf("re-encode: % x ≠ % x", buf2.Bytes(), buf.Bytes())
			return false
		}

		return true
	})
}

// TestPropertyRiscvTextRoundTripList - the "text" property for a list of
// instructions: the joined canonical text of the list assembles and decodes
// line by line into the same texts (canonical - see rvTextRoundTrip).
func TestPropertyRiscvTextRoundTripList(t *testing.T) {
	rnd := seedRnd(t)
	seq := arb.Seq(rnd, rvPropFamilies(rnd))
	ohsnap.Check(t, 100, seq, func(ins []riscv.Instr) bool {
		raw, ok := rvEncodeAll(t, ins)
		if !ok {
			return false
		}

		buf := *bytes.NewBuffer(raw)

		back, err := riscv.Parse(propAddr)(parsecbytes.Buffer(buf.Bytes()))
		if err != nil {
			t.Logf("decode: %v", err)
			return false
		}

		if len(back) != len(ins) {
			t.Logf("decoded %d of %d", len(back), len(ins))
			return false
		}

		texts := make([]string, len(back))
		for i := range back {
			texts[i] = rvText(back[i])
		}

		data, ok := rvAssemblesTo(t, strings.Join(texts, "\n"))
		if !ok {
			return false
		}

		back2, err := riscv.Parse(propAddr)(parsecbytes.Buffer(data))
		if err != nil {
			t.Logf("decode: %v", err)
			return false
		}

		if len(back2) != len(ins) {
			t.Logf("decoded %d of %d", len(back2), len(ins))
			return false
		}

		for i := range back2 {
			if rvText(back2[i]) != texts[i] {
				t.Logf("[%d] text %q ≠ %q", i, rvText(back2[i]), texts[i])
				return false
			}
		}

		return true
	})
}

// TestPropertyRiscvDecodeRobustness - arbitrary bytes do not crash the decoder.
func TestPropertyRiscvDecodeRobustness(t *testing.T) {
	rnd := seedRnd(t)
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
			ins, err := riscv.Parse(propAddr)(parsecbytes.Buffer(data))
			if err != nil {
				t.Errorf("parse %#08x: %v", w, err)
				ok = false
				return
			}

			for _, in := range ins {
				// render errors are acceptable (garbage words) - we only
				// check that no panic occurs
				odErr := in.ObjDump(disasm.DefaultViewCtx())
				_, mjErr := in.MarshalJSON()
				_, _ = odErr, mjErr
				if n := in.Len(); n != 2 && n != 4 {
					t.Logf("word %#08x: Len = %d", w, n)
					ok = false
				}
			}

			// the tail may be a truncated 32-bit form - arbitrary bytes are
			// not required to parse entirely; the requirement is not to crash
		}()
		return ok
	})
}

// TestPropertyRiscvVsObjdump - the differential: words produced by the riscv
// constructors are disassembled identically by us and by objdump.
func TestPropertyRiscvVsObjdump(t *testing.T) {
	const perFamily = 48
	const threshold = 90.0

	rnd := seedRnd(t)
	gens := rvPropFamilies(rnd)
	ins := make([]riscv.Instr, 0, perFamily*len(gens))
	for _, gen := range gens {
		for range perFamily {
			ins = append(ins, gen())
		}
	}

	raw, ok := rvEncodeAll(t, ins)
	require.True(t, ok)
	code := raw

	path := writeRiscvELF(t, code)
	out, err := objdump.Run(context.Background(), objdump.Args("ELF", 0, path))
	if err != nil {
		t.Skipf("no objdump for ELF/riscv64: %v", err)
	}

	objLines := objdump.ParseByAddr(string(out))
	require.NotEmpty(t, objLines)

	style := text.StyleFor("ELF")
	opts := disasm.NewOptions(style)
	matched, mismatched, notInOurs := 0, 0, 0
	var samples []string
	for addr, objLine := range objLines {
		off := int(addr) - propAddr
		if off < 0 || off >= len(code) {
			notInOurs++
			continue
		}

		ins, err := riscv.Parse(addr)(parsecbytes.Buffer(code[off:]))
		if err != nil {
			notInOurs++
			continue
		}

		if len(ins) == 0 {
			notInOurs++
			continue
		}

		ourLine := objdump.StripComments(objdump.Normalize(
			disasm.Line(addr, code[off:off+ins[0].Len()], ins[0], opts)))
		if ourLine == objdump.StripComments(objLine) {
			matched++
		} else {
			mismatched++
			if len(samples) < 10 {
				samples = append(samples, fmt.Sprintf(
					"0x%x\n    ours:    %s\n    objdump: %s",
					addr,
					ourLine,
					objdump.StripComments(objLine),
				))
			}
		}
	}

	total := matched + mismatched
	pct := 0.0
	if total > 0 {
		pct = float64(matched) * 100 / float64(total)
	}

	t.Logf("vs objdump: %d matched (%.2f%%), %d mismatched, %d foreign addresses",
		matched, pct, mismatched, notInOurs)
	for _, sm := range samples {
		t.Log(sm)
	}

	require.GreaterOrEqual(t, pct, threshold)
}

// writeRiscvELF - minimal ELF64 LE riscv (e_machine = EM_RISCV,
// e_flags = RVC, as in tests/examples/hello-riscv): .text @ 0x1000.
func writeRiscvELF(t *testing.T, code []byte) string {
	t.Helper()
	path, err := writeObjELF(t, code, 0xF3, 0x5)
	require.NoError(t, err)
	return path
}
