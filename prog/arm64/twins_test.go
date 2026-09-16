package arm64

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"

	arch "github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/prog"
)

// assemble - Build+Assemble with the checks every twin test needs.
func assemble(t *testing.T, p *Program) *prog.Result {
	t.Helper()

	bin, buildErrs := p.Build()
	require.Empty(t, buildErrs)

	res := bin.Assemble(0)
	require.Empty(t, res.Errs)

	return res
}

// words - assembled code as little-endian words.
func words(b []byte) []uint32 {
	out := make([]uint32, len(b)/4)
	for i := range out {
		out[i] = binary.LittleEndian.Uint32(b[i*4:])
	}

	return out
}

// enc - the word of a directly built instruction: the oracle the twins
// are byte-compared against. Panics on a construction error - the
// oracle operands are valid by construction (the pre-minted-reg
// pattern of regs.go).
func enc(i arch.Instr, err error) uint32 {
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	if _, err := i.Encode(&buf); err != nil {
		panic(err)
	}

	return binary.LittleEndian.Uint32(buf.Bytes())
}

// wreg - a W register for the byte/halfword twins.
func wreg(t *testing.T, n int) arch.Reg {
	t.Helper()

	r, err := arch.W(n)
	require.NoError(t, err)

	return r
}

func TestTwinsArithmetic(t *testing.T) {
	res := assemble(t, New().
		Adc(X0, X1, X2).
		AddImm(X3, X4, 0x123, arch.NoSh12).
		AddShift(X5, X6, X7, 5, arch.LSL).
		AddExt(X16, X29, X30, "sxtx", 0).
		AddsImm(X0, X1, 0x123, arch.LSL12).
		AddsShift(X2, X3, X4, 5, arch.LSR).
		AddsExt(X5, X6, X7, "sxtx", 1).
		SubImm(X16, X29, 0x456, arch.NoSh12).
		SubShift(X30, X0, X1, 6, arch.ASR).
		SubExt(X2, X3, X4, "sxtx", 2).
		SubsImm(XZR, X5, 0x456, arch.NoSh12). // the cmp form
		SubsShift(XZR, X6, X7, 6, arch.ASR).
		SubsExt(X0, X1, X2, "sxtx", 3).
		Madd(X3, X4, X5, XZR).  // the mul form
		Msub(X6, X7, X16, XZR). // the mneg form
		Sdiv(X29, X30, X0).
		Udiv(X1, X2, X3).
		Smulh(X4, X5, X6).
		Umulh(X7, X16, X29))

	b := arch.New()
	v12a, err := b.Imm12(0x123)
	require.NoError(t, err)
	v12b, err := b.Imm12(0x456)
	require.NoError(t, err)
	v6a, err := b.Imm6(5)
	require.NoError(t, err)
	v6b, err := b.Imm6(6)
	require.NoError(t, err)

	require.Equal(t, []uint32{
		enc(b.Adc(X0, X1, X2)),
		enc(b.AddImm(X3, X4, v12a, arch.NoSh12)),
		enc(b.AddShift(X5, X6, X7, v6a, arch.LSL)),
		enc(b.AddExt(X16, X29, X30, "sxtx", 0)),
		enc(b.AddsImm(X0, X1, v12a, arch.LSL12)),
		enc(b.AddsShift(X2, X3, X4, v6a, arch.LSR)),
		enc(b.AddsExt(X5, X6, X7, "sxtx", 1)),
		enc(b.SubImm(X16, X29, v12b, arch.NoSh12)),
		enc(b.SubShift(X30, X0, X1, v6b, arch.ASR)),
		enc(b.SubExt(X2, X3, X4, "sxtx", 2)),
		enc(b.SubsImm(XZR, X5, v12b, arch.NoSh12)),
		enc(b.SubsShift(XZR, X6, X7, v6b, arch.ASR)),
		enc(b.SubsExt(X0, X1, X2, "sxtx", 3)),
		enc(b.Madd(X3, X4, X5, XZR)),
		enc(b.Msub(X6, X7, X16, XZR)),
		enc(b.Sdiv(X29, X30, X0)),
		enc(b.Udiv(X1, X2, X3)),
		enc(b.Smulh(X4, X5, X6)),
		enc(b.Umulh(X7, X16, X29)),
	}, words(res.Code))
}

func TestTwinsLogical(t *testing.T) {
	res := assemble(t, New().
		AndImm(X0, X1, 0xff).
		AndShift(X2, X3, X4, 3, arch.LSL).
		AndsImm(X5, X6, 0x7f).
		AndsShift(X7, X16, X29, 3, arch.LSR).
		BicShift(X30, X0, X1, 4, arch.ASR).
		BicsShift(X2, X3, X4, 4, arch.ROR).
		OrrImm(X5, X6, 0xffff).
		OrrShift(X7, X16, X29, 7, arch.LSL).
		OrnShift(X30, X0, X1, 7, arch.LSR).
		EorImm(X2, X3, 0xffffffff).
		EorShift(X4, X5, X6, 9, arch.ASR).
		EonShift(X7, X16, X29, 9, arch.ROR).
		Ccmp(X0, X1, 0xf, "eq").
		Csel(X2, X3, X4, "ne").
		Csinc(X5, X6, X7, "hs").
		Csinv(X16, X29, X30, "lo").
		Csneg(X0, X1, X2, "mi"))

	b := arch.New()
	v6a, err := b.Imm6(3)
	require.NoError(t, err)
	v6b, err := b.Imm6(4)
	require.NoError(t, err)
	v6c, err := b.Imm6(7)
	require.NoError(t, err)
	v6d, err := b.Imm6(9)
	require.NoError(t, err)

	require.Equal(t, []uint32{
		enc(b.AndImm(X0, X1, 0xff)),
		enc(b.AndShift(X2, X3, X4, v6a, arch.LSL)),
		enc(b.AndsImm(X5, X6, 0x7f)),
		enc(b.AndsShift(X7, X16, X29, v6a, arch.LSR)),
		enc(b.BicShift(X30, X0, X1, v6b, arch.ASR)),
		enc(b.BicsShift(X2, X3, X4, v6b, arch.ROR)),
		enc(b.OrrImm(X5, X6, 0xffff)),
		enc(b.OrrShift(X7, X16, X29, v6c, arch.LSL)),
		enc(b.OrnShift(X30, X0, X1, v6c, arch.LSR)),
		enc(b.EorImm(X2, X3, 0xffffffff)),
		enc(b.EorShift(X4, X5, X6, v6d, arch.ASR)),
		enc(b.EonShift(X7, X16, X29, v6d, arch.ROR)),
		enc(b.Ccmp(X0, X1, 0xf, "eq")),
		enc(b.Csel(X2, X3, X4, "ne")),
		enc(b.Csinc(X5, X6, X7, "hs")),
		enc(b.Csinv(X16, X29, X30, "lo")),
		enc(b.Csneg(X0, X1, X2, "mi")),
	}, words(res.Code))
}

func TestTwinsBitfield(t *testing.T) {
	res := assemble(t, New().
		LslReg(X0, X1, X2).
		LsrReg(X3, X4, X5).
		AsrReg(X6, X7, X16).
		RorReg(X29, X30, X0).
		Bfm(X1, X2, 5, 10).
		Sbfm(X3, X4, 5, 10).
		Ubfm(X5, X6, 5, 10).
		Extr(X7, X16, X29, 12).
		Cls(X0, X1).
		Clz(X2, X3).
		Rbit(X4, X5).
		Rev(X6, X7).
		Rev16(X16, X29).
		Rev32(X30, X0))

	b := arch.New()
	v6, err := b.Imm6(12)
	require.NoError(t, err)

	require.Equal(t, []uint32{
		enc(b.LslReg(X0, X1, X2)),
		enc(b.LsrReg(X3, X4, X5)),
		enc(b.AsrReg(X6, X7, X16)),
		enc(b.RorReg(X29, X30, X0)),
		enc(b.Bfm(X1, X2, 5, 10)),
		enc(b.Sbfm(X3, X4, 5, 10)),
		enc(b.Ubfm(X5, X6, 5, 10)),
		enc(b.Extr(X7, X16, X29, v6)),
		enc(b.Cls(X0, X1)),
		enc(b.Clz(X2, X3)),
		enc(b.Rbit(X4, X5)),
		enc(b.Rev(X6, X7)),
		enc(b.Rev16(X16, X29)),
		enc(b.Rev32(X30, X0)),
	}, words(res.Code))
}

func TestTwinsControl(t *testing.T) {
	res := assemble(t, New().
		Movn(X0, 1, arch.Hw0).
		Br(X16).
		Blr(X30).
		Ret(X30).
		Nop().
		Brk(0x10).
		Mrs(X0, "CNTFRQ_EL0").
		Msr("TPIDR_EL0", X0).
		Prfm(X1))

	b := arch.New()
	v16, err := b.Imm16(1)
	require.NoError(t, err)
	v16brk, err := b.Imm16(0x10)
	require.NoError(t, err)

	require.Equal(t, []uint32{
		enc(b.Movn(X0, v16, arch.Hw0)),
		enc(b.Br(X16)),
		enc(b.Blr(X30)),
		enc(b.Ret(X30)),
		enc(b.Nop(), nil),
		enc(b.Brk(v16brk), nil),
		enc(b.Mrs(X0, "CNTFRQ_EL0")),
		enc(b.Msr("TPIDR_EL0", X0)),
		enc(b.Prfm(X1)),
	}, words(res.Code))
}

func TestTwinsMemory(t *testing.T) {
	w0 := wreg(t, 0)
	w2 := wreg(t, 2)
	w4 := wreg(t, 4)
	w1 := wreg(t, 1)
	w6 := wreg(t, 6)

	res := assemble(t, New().
		Ldr(X0, X1, 8).
		Ldrb(w2, X3, 5).
		Ldrh(w4, X5, 6).
		Ldrsb(X6, X7, 5).
		Ldrsh(X16, X29, 6).
		Ldrsw(X30, X0, 8).
		Ldur(X1, X2, -8).
		Ldurb(w2, X4, 3).
		Ldurh(w6, X6, -6).
		Ldp(X7, X16, X29, 16).
		Ldpsw(X30, X0, X1, 16).
		Str(X2, X3, 8).
		Strb(w4, X5, 5).
		Strh(w6, X7, 6).
		Stur(X16, X29, -8).
		Sturb(w2, X0, 3).
		Sturh(w4, X1, -6).
		Stp(X3, X4, X5, 16).
		Ldar(X6, X7).
		Ldarb(w0, X16).
		Ldaxr(X29, X30).
		Ldaxrb(w2, X0).
		Stlr(X1, X2).
		Stlrb(w0, X3).
		Stlxr(w1, X4, X5).
		Stlxrb(w1, w2, X6).
		Stxrb(w4, w2, X7))

	b := arch.New()
	require.Equal(t, []uint32{
		enc(b.Ldr(X0, X1, 8)),
		enc(b.Ldrb(w2, X3, 5)),
		enc(b.Ldrh(w4, X5, 6)),
		enc(b.Ldrsb(X6, X7, 5)),
		enc(b.Ldrsh(X16, X29, 6)),
		enc(b.Ldrsw(X30, X0, 8)),
		enc(b.Ldur(X1, X2, -8)),
		enc(b.Ldurb(w2, X4, 3)),
		enc(b.Ldurh(w6, X6, -6)),
		enc(b.Ldp(X7, X16, X29, 16)),
		enc(b.Ldpsw(X30, X0, X1, 16)),
		enc(b.Str(X2, X3, 8)),
		enc(b.Strb(w4, X5, 5)),
		enc(b.Strh(w6, X7, 6)),
		enc(b.Stur(X16, X29, -8)),
		enc(b.Sturb(w2, X0, 3)),
		enc(b.Sturh(w4, X1, -6)),
		enc(b.Stp(X3, X4, X5, 16)),
		enc(b.Ldar(X6, X7)),
		enc(b.Ldarb(w0, X16)),
		enc(b.Ldaxr(X29, X30)),
		enc(b.Ldaxrb(w2, X0)),
		enc(b.Stlr(X1, X2)),
		enc(b.Stlrb(w0, X3)),
		enc(b.Stlxr(w1, X4, X5)),
		enc(b.Stlxrb(w1, w2, X6)),
		enc(b.Stxrb(w4, w2, X7)),
	}, words(res.Code))
}

func TestTwinsLabelDirected(t *testing.T) {
	// base 0xff0: the adrp sits in page 0, the page label at +0x2c
	// lands in page 1; tbz covers a forward and a backward target
	w1 := wreg(t, 1)
	p := New().
		Adrp(X0, "page").        // @0x00: page delta +1
		Label("top").            // @0x04
		Tbz(w1, 3, "fwd").       // @0x04 → fwd @0x28: +36
		Nop().Nop().Nop().Nop(). // @0x08..
		Nop().Nop().Nop().Nop(). // ..@0x24
		Label("fwd").            // @0x28
		Tbz(X2, 40, "top").      // @0x28 → top @0x04: -36 (bit 40 → x)
		Label("page")            // @0x2c: abs 0x101c, page 1

	bin, buildErrs := p.Build()
	require.Empty(t, buildErrs)

	res := bin.Assemble(0xff0)
	require.Empty(t, res.Errs)

	b := arch.New()
	require.Equal(t, []uint32{
		enc(b.Adrp(X0, 1)),
		enc(b.Tbz(w1, 3, 36)),
		enc(b.Nop(), nil),
		enc(b.Nop(), nil),
		enc(b.Nop(), nil),
		enc(b.Nop(), nil),
		enc(b.Nop(), nil),
		enc(b.Nop(), nil),
		enc(b.Nop(), nil),
		enc(b.Nop(), nil),
		enc(b.Tbz(X2, 40, -36)),
	}, words(res.Code))
}

func TestTwinsErrors(t *testing.T) {
	// imm out of range: the deferred error surfaces at Build
	_, buildErrs := New().AddImm(X0, X1, 0x1000, arch.NoSh12).Build()
	require.Len(t, buildErrs, 1)
	require.Contains(t, buildErrs[0].Error(), "add")

	// misaligned load offset: the ctor error surfaces at Build
	_, buildErrs = New().Ldr(X0, X1, 5).Build()
	require.Len(t, buildErrs, 1)
	require.Contains(t, buildErrs[0].Error(), "ldr")

	// label-directed lines build fine and surface their errors at
	// Assemble: the tbz width rule (bit 40 needs an x register)
	bin, buildErrs := New().Tbz(wreg(t, 1), 40, "l").Label("l").Build()
	require.Empty(t, buildErrs)
	errs := bin.Assemble(0).Errs
	require.Len(t, errs, 1)
	require.Contains(t, errs[0].Error(), "tbz")

	// ...and the undefined target
	bin, buildErrs = New().Tbz(X2, 5, "nowhere").Build()
	require.Empty(t, buildErrs)
	errs = bin.Assemble(0).Errs
	require.Len(t, errs, 1)
	require.Contains(t, errs[0].Error(), "nowhere")
}
