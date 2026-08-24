package assembly_test

// Property tests (oh-snap), arm64: round trip of single instructions and
// compositions, decoder robustness on arbitrary words, differential against
// the real objdump. Generators are arb/arm64 on top of arch/arm64
// constructors; the oracle is assemble(ObjDump(instr)) == instr.
// The common part of the suite is property_test.go.

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
	a64 "github.com/okneniz/assembly/arb/arm64"
	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/asm/arm64/alias"
	"github.com/okneniz/assembly/disasm"
	"github.com/okneniz/assembly/tests/cmd/objdump"
	"github.com/okneniz/assembly/text"
)

// propText - normalized ObjDump text of an instruction.
func propText(in arm64.Instr) string {
	return objdump.StripComments(objdump.Normalize(in.ObjDump(disasm.DefaultViewCtx())))
}

// bytesOf - instruction bytes (Encode encoding).
func bytesOf(t *testing.T, in arm64.Instr) ([]byte, bool) {
	t.Helper()
	var buf bytes.Buffer
	if _, err := in.Encode(&buf, propAddr); err != nil {
		t.Logf("%s: Encode: %v", in.ObjDump(disasm.DefaultViewCtx()), err)
		return nil, false
	}

	return buf.Bytes(), true
}

// a64Enc - context of the "bytes" property: a fixed address.
type a64Enc struct {
	addr uint64
}

// Addr - the address of the context (RoundTrip decoder).
func (c a64Enc) Addr() uint64 {
	return c.addr
}

// a64EncodeAll - encodes a list sequentially starting at propAddr (the address
// of each is propAddr + bytes written; no symbols).
func a64EncodeAll(t *testing.T, ins []arm64.Instr) ([]byte, bool) {
	t.Helper()
	var buf bytes.Buffer
	addr := propAddr
	for _, in := range ins {
		n, err := in.Encode(&buf, uint64(addr))
		if err != nil {
			t.Logf("encode: %v", err)
			return nil, false
		}

		addr += int(n)
	}

	return buf.Bytes(), true
}

// assemblesTo - bytes from assembling the text (false on assembly error).
func assemblesTo(t *testing.T, src string) ([]byte, bool) {
	t.Helper()
	res, errs := alias.Assemble(src, propAddr)
	if len(errs) != 0 {
		t.Logf("%q: assemble: %v", src, errs)
		return nil, false
	}

	return res.Sections[0].Data, true
}

// propBytesRoundTrip - the "bytes" property: the law enc∘dec∘enc == enc
// (RoundTrip) in the encoding context propAddr without symbols - bytes are
// stable after the round trip.
func propBytesRoundTrip(t *testing.T, in arm64.Instr) bool {
	t.Helper()
	return RoundTrip[a64Enc, arm64.Instr, []byte](
		a64Enc{addr: propAddr},
		func(_ a64Enc, x arm64.Instr) ([]byte, bool) {
			return bytesOf(t, x)
		},
		func(ctx a64Enc, b []byte) (arm64.Instr, bool) {
			back, err := arm64.Parse(ctx.addr)(parsecbytes.Buffer(b))
			if err != nil {
				t.Logf("decode: %v", err)
				return nil, false
			}

			if len(back) != 1 {
				t.Logf("%#08x: decode = %d instr", binary.LittleEndian.Uint32(b), len(back))
				return nil, false
			}

			return back[0], true
		},
		bytes.Equal,
	)(in)
}

// propTextRoundTrip - the "text" property: the RoundTrip law in the
// DefaultViewCtx context - the canonical text of an instruction assembles and
// decodes back into the same text.
func propTextRoundTrip(t *testing.T, in arm64.Instr) bool {
	t.Helper()
	return RoundTrip[disasm.ViewCtx, arm64.Instr, string](
		disasm.DefaultViewCtx(),
		func(_ disasm.ViewCtx, x arm64.Instr) (string, bool) {
			return propText(x), true
		},
		func(_ disasm.ViewCtx, src string) (arm64.Instr, bool) {
			data, ok := assemblesTo(t, src)
			if !ok {
				return nil, false
			}

			back, err := arm64.Parse(propAddr)(parsecbytes.Buffer(data))
			if err != nil {
				t.Logf("decode: %v", err)
				return nil, false
			}

			if len(back) != 1 {
				t.Logf("%q: decode = %d instr", src, len(back))
				return nil, false
			}

			return back[0], true
		},
		func(a, b string) bool { return a == b },
	)(in)
}

// instrParam - parameters of any family (arb generators return exactly these).
type instrParam interface {
	Instr() arm64.Instr
}

// propFamilyEntry - one family in the TestPropertySingleInstrRoundTrip table.
type propFamilyEntry struct {
	name string
	run  func(t *testing.T, rnd *mrnd.Rand)
}

// newPropFamily - a family entry: closes over the generic instantiation of
// the properties. The parameter types of families differ (RetParams,
// SvcParams, ...), while the generic ohsnap.Arbitrary interface is invariant -
// a common Arbitrary[instrParam] cannot be assembled for the table, so the
// closure is written here, once. Each property is a separate named subtest
// (bytes / text); the check runs over the parameters for the sake of
// shrinking (ohsnap.Map onto the instruction would have truncated the shrink).
func newPropFamily[P instrParam](
	name string,
	mk func(rnd *mrnd.Rand) ohsnap.Arbitrary[P],
) propFamilyEntry {
	return propFamilyEntry{
		name: name,
		run: func(t *testing.T, rnd *mrnd.Rand) {
			t.Helper()

			t.Run("bytes", func(t *testing.T) {
				ohsnap.Check(t, 100000, mk(rnd), func(p P) bool {
					return propBytesRoundTrip(t, p.Instr())
				})
			})

			t.Run("text", func(t *testing.T) {
				ohsnap.Check(t, 100000, mk(rnd), func(p P) bool {
					return propTextRoundTrip(t, p.Instr())
				})
			})
		},
	}
}

// TestPropertySingleInstrRoundTrip - round trip of each family separately.
func TestPropertySingleInstrRoundTrip(t *testing.T) {
	families := []propFamilyEntry{
		newPropFamily("Ret", a64.Ret),
		newPropFamily("Svc", a64.Svc),
		newPropFamily("Brk", a64.Brk),
		newPropFamily("Movz", a64.Movz),
		newPropFamily("Movk", a64.Movk),
		newPropFamily("AddImm", a64.AddImm),
		newPropFamily("SubImm", a64.SubImm),
		newPropFamily("AddShift", a64.AddShift),
		newPropFamily("SubShift", a64.SubShift),
		newPropFamily("Ldr", a64.Ldr),
		newPropFamily("Str", a64.Str),
	}
	for _, f := range families {
		t.Run(f.name, func(t *testing.T) {
			f.run(t, seedRnd(t))
		})
	}
}

// newAliasFamily - like newPropFamily, but the parameters are first pinned
// (alias condition: a specific register/immediate under which the instruction
// is written as an alias). Combinations that are invalid after pinning (the
// constructor refused) are outside the property.
func newAliasFamily[P instrParam](
	name string,
	mk func(rnd *mrnd.Rand) ohsnap.Arbitrary[P],
	pin func(P) P,
) propFamilyEntry {
	return propFamilyEntry{
		name: name,
		run: func(t *testing.T, rnd *mrnd.Rand) {
			t.Helper()
			pinned := func(p P) (arm64.Instr, bool) {
				q := pin(p)
				if q.Instr() == nil {
					return nil, false
				}

				return q.Instr(), true
			}

			t.Run("bytes", func(t *testing.T) {
				ohsnap.Check(t, 100000, mk(rnd), func(p P) bool {
					if in, ok := pinned(p); ok {
						return propBytesRoundTrip(t, in)
					}

					return true
				})
			})

			t.Run("text", func(t *testing.T) {
				ohsnap.Check(t, 100000, mk(rnd), func(p P) bool {
					if in, ok := pinned(p); ok {
						return propTextRoundTrip(t, in)
					}

					return true
				})
			})
		},
	}
}

// TestPropertyAliasRoundTrip - the "alias" property: an instruction written
// as an alias survives the cycle struct → alias → struct without loss - the
// alias text assembles and decodes back into the same instruction (each entry
// is its own subtest). mov is covered by the Movz family: a struct movz always
// renders as "mov #imm".
func TestPropertyAliasRoundTrip(t *testing.T) {
	families := []propFamilyEntry{
		newAliasFamily("neg", a64.SubShift, func(p a64.SubShiftParams) a64.SubShiftParams {
			p.Rn = arm64.XZR
			return p
		}),
	}
	for _, f := range families {
		t.Run(f.name, func(t *testing.T) {
			f.run(t, seedRnd(t))
		})
	}
}

// instrOf - a sampler function over a family generator: the generator is
// created once and closed over; each sampler call is just Generate.
func instrOf[P instrParam](a ohsnap.Arbitrary[P]) func() arm64.Instr {
	return func() arm64.Instr {
		return a.Generate().Instr()
	}
}

// propFamilies - family generators as sources of instructions for
// composition and the differential.
func propFamilies(rnd *mrnd.Rand) []func() arm64.Instr {
	return []func() arm64.Instr{
		instrOf(a64.Ret(rnd)),
		instrOf(a64.Svc(rnd)),
		instrOf(a64.Brk(rnd)),
		instrOf(a64.Movz(rnd)),
		instrOf(a64.Movk(rnd)),
		instrOf(a64.AddImm(rnd)),
		instrOf(a64.SubImm(rnd)),
		instrOf(a64.AddShift(rnd)),
		instrOf(a64.SubShift(rnd)),
		instrOf(a64.Ldr(rnd)),
		instrOf(a64.Str(rnd)),
	}
}

// TestPropertyBytesRoundTripList - the "bytes" property for a list of
// instructions: the list is encoded, decoded line by line, and encoded again
// into the same bytes (struct → bytes → struct, without loss).
func TestPropertyBytesRoundTripList(t *testing.T) {
	rnd := seedRnd(t)
	seq := arb.Seq(rnd, propFamilies(rnd))

	ohsnap.Check(t, 100000, seq, func(ins []arm64.Instr) bool {
		raw, ok := a64EncodeAll(t, ins)
		if !ok {
			return false
		}

		buf := *bytes.NewBuffer(raw)

		back, err := arm64.Parse(propAddr)(parsecbytes.Buffer(buf.Bytes()))
		if err != nil {
			t.Logf("decode: %v", err)
			return false
		}

		if len(back) != len(ins) {
			t.Logf("decoded %d of %d", len(back), len(ins))
			return false
		}

		raw2, ok := a64EncodeAll(t, back)
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

// TestPropertyTextRoundTripList - the "text" property for a list of
// instructions: the joined objdump text of the list assembles and decodes
// line by line into the same texts (struct → objdump → struct, without loss).
func TestPropertyTextRoundTripList(t *testing.T) {
	rnd := seedRnd(t)
	seq := arb.Seq(rnd, propFamilies(rnd))

	ohsnap.Check(t, 100000, seq, func(ins []arm64.Instr) bool {
		texts := make([]string, len(ins))
		for i, in := range ins {
			texts[i] = propText(in)
		}

		data, ok := assemblesTo(t, strings.Join(texts, "\n"))
		if !ok {
			return false
		}

		back, err := arm64.Parse(propAddr)(parsecbytes.Buffer(data))
		if err != nil {
			t.Logf("decode: %v", err)
			return false
		}

		if len(back) != len(ins) {
			t.Logf("decoded %d of %d", len(back), len(ins))
			return false
		}

		for i := range back {
			if propText(back[i]) != texts[i] {
				t.Logf("[%d] text %q ≠ %q", i, propText(back[i]), texts[i])
				return false
			}
		}

		return true
	})
}

// TestPropertyDecodeRobustness - an arbitrary word does not crash the
// decoder: Parse + ObjDump + MarshalJSON without panics, one instruction of
// length 4.
func TestPropertyDecodeRobustness(t *testing.T) {
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
			ins, err := arm64.Parse(propAddr)(parsecbytes.Buffer(data))
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
			}

			ok = len(ins) == 1 && ins[0].Len() == 4
		}()
		return ok
	})
}

// TestPropertyArm64VsObjdump - the differential: words produced by the
// constructors are disassembled identically by us and by objdump (by
// addresses, normalized strings; threshold same as in the objdump gate).
func TestPropertyArm64VsObjdump(t *testing.T) {
	const perFamily = 48
	const threshold = 90.0

	rnd := seedRnd(t)
	gens := propFamilies(rnd)
	ins := make([]arm64.Instr, 0, perFamily*len(gens))
	for _, gen := range gens {
		for range perFamily {
			ins = append(ins, gen())
		}
	}

	code, ok := a64EncodeAll(t, ins)
	require.True(t, ok)

	elfPath, err := writeObjELF(t, code, 0xB7, 0)
	require.NoError(t, err)
	path := elfPath
	out, err := objdump.Run(context.Background(), objdump.Args("ELF", 0, path))
	if err != nil {
		t.Skipf("no objdump for ELF/arm64: %v", err)
	}

	objLines := objdump.ParseByAddr(string(out))
	require.NotEmpty(t, objLines)

	// Our lines use the same columnar form (ELF - hex word).
	style := text.StyleFor("ELF")
	opts := disasm.NewOptions(style)
	matched, mismatched, notInOurs := 0, 0, 0
	var samples []string
	for addr, objLine := range objLines {
		off := int(addr) - propAddr
		if off < 0 || off+4 > len(code) {
			notInOurs++
			continue
		}

		ours, err := arm64.Parse(addr)(parsecbytes.Buffer(code[off : off+4]))
		if err != nil {
			notInOurs++
			continue
		}

		if len(ours) != 1 {
			notInOurs++
			continue
		}

		ourLine := objdump.StripComments(objdump.Normalize(
			disasm.Line(addr, code[off:off+4], ours[0], opts)))
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
	for _, s := range samples {
		t.Log(s)
	}

	require.GreaterOrEqual(t, pct, threshold)
}
