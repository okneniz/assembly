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
	"github.com/okneniz/parsec"
	parsecbytes "github.com/okneniz/parsec/bytes"
	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/arb"
	a64 "github.com/okneniz/assembly/arb/arm64"
	acmp "github.com/okneniz/assembly/arb/arm64/alias"
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

// propTextAt - propText in the context of an explicit base address.
func propTextAt(in arm64.Instr, base uint64) string {
	return objdump.StripComments(objdump.Normalize(in.ObjDump(disasm.ViewCtxAt(base))))
}

// bytesOf - instruction bytes (Encode encoding).
func bytesOf(t *testing.T, in arm64.Instr) ([]byte, bool) {
	t.Helper()
	var buf bytes.Buffer
	if _, err := in.Encode(&buf); err != nil {
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
		n, err := in.Encode(&buf)
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

// assemblesToAt - assemblesTo at an explicit base: the pc-relative texts
// carry absolute targets, so the base must match the render context.
func assemblesToAt(t *testing.T, src string, base uint64) ([]byte, bool) {
	t.Helper()
	res, errs := alias.Assemble(src, base)
	if len(errs) != 0 {
		t.Logf("%q: assemble@%#x: %v", src, base, errs)
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
			back, err := arm64.MakeDecoder()(parsec.Stateless{}, parsecbytes.Buffer(b))
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
	return propTextRoundTripAt(t, in, 0)
}

// propTextRoundTripAt - the text law at an explicit base: the text renders
// from ViewCtxAt(base) (pc-relative operands print absolute targets of that
// base) and assembles back at the same base. Base 0 is the DefaultViewCtx
// render; propAddr is the historic assembly base - both must hold.
func propTextRoundTripAt(t *testing.T, in arm64.Instr, base uint64) bool {
	t.Helper()
	return RoundTrip[disasm.ViewCtx, arm64.Instr, string](
		disasm.ViewCtxAt(base),
		func(_ disasm.ViewCtx, x arm64.Instr) (string, bool) {
			return propTextAt(x, base), true
		},
		func(_ disasm.ViewCtx, src string) (arm64.Instr, bool) {
			data, ok := assemblesToAt(t, src, base)
			if !ok {
				return nil, false
			}

			back, err := arm64.MakeDecoder()(parsec.Stateless{}, parsecbytes.Buffer(data))
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

// textShrinkBudget - the cap on property calls spent shrinking a failing
// text case: every candidate costs a full Assemble+Parse round trip, an
// unbounded minimization of a broken encoder would stall the suite. A
// descent over an instruction's shrink axes takes hundreds of calls, so
// 1000 keeps a full minimization and cuts the pathological ones.
const textShrinkBudget = 1000

// checkProgressEvery - a progress message every N passed cases: a long
// property run confirms being alive in the log (visible with go test -v).
const checkProgressEvery = 10000

// checkOpts - the common options of the suite's property checks: progress
// lines and the accepted shrinking steps in the log.
func checkOpts(budget int) ohsnap.CheckOptions {
	return ohsnap.CheckOptions{
		Budget:         budget,
		ProgressEvery:  checkProgressEvery,
		LogShrinkSteps: true,
	}
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
				ohsnap.CheckWith(t, 100000, mk(rnd), func(p P) bool {
					return propBytesRoundTrip(t, p.Instr())
				}, checkOpts(0))
			})

			t.Run("text", func(t *testing.T) {
				// both bases: 0 (the DefaultViewCtx render) and propAddr -
				// the pc-relative texts carry absolute targets, the base
				// shift must not break the law anywhere.
				for _, base := range []uint64{0, propAddr} {
					ohsnap.CheckWith(t, 100000, mk(rnd), func(p P) bool {
						return propTextRoundTripAt(t, p.Instr(), base)
					}, checkOpts(textShrinkBudget))
				}
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
		newPropFamily("Nop", a64.Nop),
		newPropFamily("Isb", a64.Isb),
		newPropFamily("Dsb", a64.Dsb),
		newPropFamily("Dmb", a64.Dmb),
		newPropFamily("Smc", a64.Smc),
		newPropFamily("Br", a64.Br),
		newPropFamily("Blr", a64.Blr),
		newPropFamily("Movn", a64.Movn),
		newPropFamily("Adc", a64.Adc),
		newPropFamily("Smulh", a64.Smulh),
		newPropFamily("Umulh", a64.Umulh),
		newPropFamily("Rev", a64.Rev),
		newPropFamily("Rev16", a64.Rev16),
		newPropFamily("Rev32", a64.Rev32),
		newPropFamily("Cls", a64.Cls),
		newPropFamily("Clz", a64.Clz),
		newPropFamily("Rbit", a64.Rbit),
		newPropFamily("Sdiv", a64.Sdiv),
		newPropFamily("Udiv", a64.Udiv),
		newPropFamily("LslReg", a64.LslReg),
		newPropFamily("LsrReg", a64.LsrReg),
		newPropFamily("AsrReg", a64.AsrReg),
		newPropFamily("RorReg", a64.RorReg),
		newPropFamily("Mrs", a64.Mrs),
		newPropFamily("Msr", a64.Msr),
		newPropFamily("B", a64.B),
		newPropFamily("Bl", a64.Bl),
		newPropFamily("Bcond", a64.Bcond),
		newPropFamily("Cbz", a64.Cbz),
		newPropFamily("Cbnz", a64.Cbnz),
		newPropFamily("Tbz", a64.Tbz),
		newPropFamily("Adr", a64.Adr),
		newPropFamily("Adrp", a64.Adrp),
		newPropFamily("AddsImm", a64.AddsImm),
		newPropFamily("AddsShift", a64.AddsShift),
		newPropFamily("AddExt", a64.AddExt),
		newPropFamily("AddsExt", a64.AddsExt),
		newPropFamily("SubExt", a64.SubExt),
		newPropFamily("SubsExt", a64.SubsExt),
		newPropFamily("AndShift", a64.AndShift),
		newPropFamily("AndsShift", a64.AndsShift),
		newPropFamily("OrrShift", a64.OrrShift),
		newPropFamily("EorShift", a64.EorShift),
		newPropFamily("OrnShift", a64.OrnShift),
		newPropFamily("EonShift", a64.EonShift),
		newPropFamily("BicShift", a64.BicShift),
		newPropFamily("BicsShift", a64.BicsShift),
		newPropFamily("AndImm", a64.AndImm),
		newPropFamily("OrrImm", a64.OrrImm),
		newPropFamily("EorImm", a64.EorImm),
		newPropFamily("AndsImm", a64.AndsImm),
		newPropFamily("Bfm", a64.Bfm),
		newPropFamily("Sbfm", a64.Sbfm),
		newPropFamily("Ubfm", a64.Ubfm),
		newPropFamily("Extr", a64.Extr),
		newPropFamily("Ldrb", a64.Ldrb),
		newPropFamily("Strb", a64.Strb),
		newPropFamily("Ldrh", a64.Ldrh),
		newPropFamily("Strh", a64.Strh),
		newPropFamily("Ldrsb", a64.Ldrsb),
		newPropFamily("Ldrsh", a64.Ldrsh),
		newPropFamily("Ldrsw", a64.Ldrsw),
		newPropFamily("Ldur", a64.Ldur),
		newPropFamily("Ldurb", a64.Ldurb),
		newPropFamily("Ldurh", a64.Ldurh),
		newPropFamily("Stur", a64.Stur),
		newPropFamily("Sturb", a64.Sturb),
		newPropFamily("Sturh", a64.Sturh),
		newPropFamily("LdrF", a64.LdrF),
		newPropFamily("StrF", a64.StrF),
		newPropFamily("Ldp", a64.Ldp),
		newPropFamily("Stp", a64.Stp),
		newPropFamily("Ldpsw", a64.Ldpsw),
		newPropFamily("Ldar", a64.Ldar),
		newPropFamily("Ldaxr", a64.Ldaxr),
		newPropFamily("Stlr", a64.Stlr),
		newPropFamily("Ldarb", a64.Ldarb),
		newPropFamily("Ldaxrb", a64.Ldaxrb),
		newPropFamily("Stlrb", a64.Stlrb),
		newPropFamily("Stlxr", a64.Stlxr),
		newPropFamily("Stxrb", a64.Stxrb),
		newPropFamily("Stlxrb", a64.Stlxrb),
		newPropFamily("Prfm", a64.Prfm),
		newPropFamily("Fadd", a64.Fadd),
		newPropFamily("Fsub", a64.Fsub),
		newPropFamily("Fmul", a64.Fmul),
		newPropFamily("Fdiv", a64.Fdiv),
		newPropFamily("Fmax", a64.Fmax),
		newPropFamily("Fmin", a64.Fmin),
		newPropFamily("Fcmp", a64.Fcmp),
		newPropFamily("FcmpZero", a64.FcmpZero),
		newPropFamily("Fneg", a64.Fneg),
		newPropFamily("Fmov", a64.Fmov),
		newPropFamily("Fcvt", a64.Fcvt),
		newPropFamily("Fmadd", a64.Fmadd),
		newPropFamily("Fnmsub", a64.Fnmsub),
		newPropFamily("FmovFromGpr", a64.FmovFromGpr),
		newPropFamily("FmovToGpr", a64.FmovToGpr),
		newPropFamily("Fcvtzs", a64.Fcvtzs),
		newPropFamily("Fcvtzu", a64.Fcvtzu),
		newPropFamily("Scvtf", a64.Scvtf),
		newPropFamily("Ucvtf", a64.Ucvtf),
		newPropFamily("FmovImm", a64.FmovImm),
	}
	for _, f := range families {
		t.Run(f.name, func(t *testing.T) {
			f.run(t, seedRnd(t))
		})
	}
}

// TestPropertyAliasRoundTrip - the "alias" property: the alias families
// of arb/arm64/alias - the parameters render the alias text, it assembles
// into the base instruction, and the round trips hold (each entry is its
// own subtest). mov is covered both ways: the Movz family renders the
// canonical mov #imm, the alias family walks the alias operand space.
func TestPropertyAliasRoundTrip(t *testing.T) {
	families := []propFamilyEntry{
		newPropFamily("Cmp", acmp.Cmp),
		newPropFamily("Cmn", acmp.Cmn),
		newPropFamily("Neg", acmp.Neg),
		newPropFamily("Negs", acmp.Negs),
		newPropFamily("Tst", acmp.Tst),
		newPropFamily("Mvn", acmp.Mvn),
		newPropFamily("Mov", acmp.Mov),
		newPropFamily("Mul", acmp.Mul),
		newPropFamily("Mneg", acmp.Mneg),
		newPropFamily("Cset", acmp.Cset),
		newPropFamily("Csetm", acmp.Csetm),
		newPropFamily("Cinc", acmp.Cinc),
		newPropFamily("Cinv", acmp.Cinv),
		newPropFamily("Cneg", acmp.Cneg),
		newPropFamily("Sxtb", acmp.Sxtb),
		newPropFamily("Sxth", acmp.Sxth),
		newPropFamily("Sxtw", acmp.Sxtw),
		newPropFamily("Ubfiz", acmp.Ubfiz),
		newPropFamily("Ubfx", acmp.Ubfx),
		newPropFamily("Sbfiz", acmp.Sbfiz),
		newPropFamily("Sbfx", acmp.Sbfx),
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
		return ohsnap.First(a.Generate()).Instr()
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
		instrOf(a64.Nop(rnd)),
		instrOf(a64.Isb(rnd)),
		instrOf(a64.Dsb(rnd)),
		instrOf(a64.Dmb(rnd)),
		instrOf(a64.Smc(rnd)),
		instrOf(a64.Br(rnd)),
		instrOf(a64.Blr(rnd)),
		instrOf(a64.Movn(rnd)),
		instrOf(a64.Adc(rnd)),
		instrOf(a64.Smulh(rnd)),
		instrOf(a64.Umulh(rnd)),
		instrOf(a64.Rev(rnd)),
		instrOf(a64.Rev16(rnd)),
		instrOf(a64.Rev32(rnd)),
		instrOf(a64.Cls(rnd)),
		instrOf(a64.Clz(rnd)),
		instrOf(a64.Rbit(rnd)),
		instrOf(a64.Sdiv(rnd)),
		instrOf(a64.Udiv(rnd)),
		instrOf(a64.LslReg(rnd)),
		instrOf(a64.LsrReg(rnd)),
		instrOf(a64.AsrReg(rnd)),
		instrOf(a64.RorReg(rnd)),
		instrOf(a64.Mrs(rnd)),
		instrOf(a64.Msr(rnd)),
		instrOf(a64.B(rnd)),
		instrOf(a64.Bl(rnd)),
		instrOf(a64.Bcond(rnd)),
		instrOf(a64.Cbz(rnd)),
		instrOf(a64.Cbnz(rnd)),
		instrOf(a64.Tbz(rnd)),
		instrOf(a64.Adr(rnd)),
		instrOf(a64.Adrp(rnd)),
		instrOf(a64.AddsImm(rnd)),
		instrOf(a64.AddsShift(rnd)),
		instrOf(a64.AddExt(rnd)),
		instrOf(a64.AddsExt(rnd)),
		instrOf(a64.SubExt(rnd)),
		instrOf(a64.SubsExt(rnd)),
		instrOf(a64.AndShift(rnd)),
		instrOf(a64.AndsShift(rnd)),
		instrOf(a64.OrrShift(rnd)),
		instrOf(a64.EorShift(rnd)),
		instrOf(a64.OrnShift(rnd)),
		instrOf(a64.EonShift(rnd)),
		instrOf(a64.BicShift(rnd)),
		instrOf(a64.BicsShift(rnd)),
		instrOf(a64.AndImm(rnd)),
		instrOf(a64.OrrImm(rnd)),
		instrOf(a64.EorImm(rnd)),
		instrOf(a64.AndsImm(rnd)),
		instrOf(a64.Bfm(rnd)),
		instrOf(a64.Sbfm(rnd)),
		instrOf(a64.Ubfm(rnd)),
		instrOf(a64.Extr(rnd)),
		instrOf(a64.Ldrb(rnd)),
		instrOf(a64.Strb(rnd)),
		instrOf(a64.Ldrh(rnd)),
		instrOf(a64.Strh(rnd)),
		instrOf(a64.Ldrsb(rnd)),
		instrOf(a64.Ldrsh(rnd)),
		instrOf(a64.Ldrsw(rnd)),
		instrOf(a64.Ldur(rnd)),
		instrOf(a64.Ldurb(rnd)),
		instrOf(a64.Ldurh(rnd)),
		instrOf(a64.Stur(rnd)),
		instrOf(a64.Sturb(rnd)),
		instrOf(a64.Sturh(rnd)),
		instrOf(a64.LdrF(rnd)),
		instrOf(a64.StrF(rnd)),
		instrOf(a64.Ldp(rnd)),
		instrOf(a64.Stp(rnd)),
		instrOf(a64.Ldpsw(rnd)),
		instrOf(a64.Ldar(rnd)),
		instrOf(a64.Ldaxr(rnd)),
		instrOf(a64.Stlr(rnd)),
		instrOf(a64.Ldarb(rnd)),
		instrOf(a64.Ldaxrb(rnd)),
		instrOf(a64.Stlrb(rnd)),
		instrOf(a64.Stlxr(rnd)),
		instrOf(a64.Stxrb(rnd)),
		instrOf(a64.Stlxrb(rnd)),
		instrOf(a64.Prfm(rnd)),
		instrOf(a64.Fadd(rnd)),
		instrOf(a64.Fsub(rnd)),
		instrOf(a64.Fmul(rnd)),
		instrOf(a64.Fdiv(rnd)),
		instrOf(a64.Fmax(rnd)),
		instrOf(a64.Fmin(rnd)),
		instrOf(a64.Fcmp(rnd)),
		instrOf(a64.FcmpZero(rnd)),
		instrOf(a64.Fneg(rnd)),
		instrOf(a64.Fmov(rnd)),
		instrOf(a64.Fcvt(rnd)),
		instrOf(a64.Fmadd(rnd)),
		instrOf(a64.Fnmsub(rnd)),
		instrOf(a64.FmovFromGpr(rnd)),
		instrOf(a64.FmovToGpr(rnd)),
		instrOf(a64.Fcvtzs(rnd)),
		instrOf(a64.Fcvtzu(rnd)),
		instrOf(a64.Scvtf(rnd)),
		instrOf(a64.Ucvtf(rnd)),
		instrOf(a64.FmovImm(rnd)),
	}
}

// TestPropertyBytesRoundTripList - the "bytes" property for a list of
// instructions: the list is encoded, decoded line by line, and encoded again
// into the same bytes (struct → bytes → struct, without loss).
func TestPropertyBytesRoundTripList(t *testing.T) {
	rnd := seedRnd(t)
	seq := arb.Seq(rnd, propFamilies(rnd))

	ohsnap.CheckWith(t, 100000, seq, func(ins []arm64.Instr) bool {
		raw, ok := a64EncodeAll(t, ins)
		if !ok {
			return false
		}

		buf := *bytes.NewBuffer(raw)

		back, err := arm64.MakeDecoder()(parsec.Stateless{}, parsecbytes.Buffer(buf.Bytes()))
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
	}, checkOpts(0))
}

// TestPropertyTextRoundTripList - the "text" property for a list of
// instructions: the joined objdump text of the list assembles and decodes
// line by line into the same texts (struct → objdump → struct, without loss).
func TestPropertyTextRoundTripList(t *testing.T) {
	rnd := seedRnd(t)
	seq := arb.Seq(rnd, propFamilies(rnd))

	ohsnap.CheckWith(t, 100000, seq, func(ins []arm64.Instr) bool {
		texts := make([]string, len(ins))
		for i, in := range ins {
			// the text renders at the instruction's own address in the
			// list (pc-relative targets print absolute) - the joined text
			// assembles at propAddr, the same base the render used.
			texts[i] = propTextAt(in, propAddr+uint64(4*i))
		}

		data, ok := assemblesTo(t, strings.Join(texts, "\n"))
		if !ok {
			return false
		}

		back, err := arm64.MakeDecoder()(parsec.Stateless{}, parsecbytes.Buffer(data))
		if err != nil {
			t.Logf("decode: %v", err)
			return false
		}

		if len(back) != len(ins) {
			t.Logf("decoded %d of %d", len(back), len(ins))
			return false
		}

		for i := range back {
			if propTextAt(back[i], propAddr+uint64(4*i)) != texts[i] {
				t.Logf("[%d] text %q ≠ %q", i, propTextAt(back[i], propAddr+uint64(4*i)), texts[i])
				return false
			}
		}

		return true
	}, checkOpts(0))
}

// TestPropertyDecodeRobustness - an arbitrary word does not crash the
// decoder: Parse + ObjDump without panics, one instruction of
// length 4.
func TestPropertyDecodeRobustness(t *testing.T) {
	rnd := seedRnd(t)

	ohsnap.CheckWith(t, 100000, arb.Word(rnd), func(w uint32) bool {
		ok := true
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("panic at %#08x: %v", w, r)
					ok = false
				}
			}()
			data := binary.LittleEndian.AppendUint32(nil, w)
			ins, err := arm64.MakeDecoder()(parsec.Stateless{}, parsecbytes.Buffer(data))
			if err != nil {
				t.Errorf("parse %#08x: %v", w, err)
				ok = false
				return
			}

			for _, in := range ins {
				// render errors are acceptable (garbage words) - we only
				// check that no panic occurs
				_ = in.ObjDump(disasm.DefaultViewCtx())
			}

			ok = len(ins) == 1 && ins[0].Len() == 4
		}()
		return ok
	}, checkOpts(0))
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

		ours, err := arm64.MakeDecoder()(parsec.Stateless{}, parsecbytes.Buffer(code[off:off+4]))
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
