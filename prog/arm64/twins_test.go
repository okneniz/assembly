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

func TestTwinsFP(t *testing.T) {
	w0 := wreg(t, 0)
	w1 := wreg(t, 1)

	res := assemble(t, New().
		Fadd(D0, D1, D2).
		Fadd(S0, S1, S2).
		Fsub(D3, D4, D5).
		Fmul(S3, S4, S5).
		Fdiv(D6, D7, D0).
		Fmax(S6, S7, S0).
		Fmin(D1, D2, D3).
		Fneg(D4, D5).
		Fneg(S4, S5).
		Fmov(D6, D7).
		Fmov(S6, S7).
		FmovFromGpr(D0, X0).
		FmovFromGpr(S1, w1).
		FmovFromGpr(D2, XZR). // the FP zero idiom
		FmovToGpr(X3, D5).
		FmovToGpr(w0, S7).
		FmovImm(D0, 0.5).
		FmovImm(S0, 1.5).
		Fcvt(S0, D0).
		Fcvt(D1, S1).
		Scvtf(D2, w1).
		Scvtf(D3, X2).
		Scvtf(S2, w1).
		Scvtf(S3, X3).
		Ucvtf(D4, X4).
		Ucvtf(S4, w0).
		Fcvtzs(X5, D5).
		Fcvtzs(w1, D6).
		Fcvtzs(X6, S5).
		Fcvtzs(w0, S6).
		Fcvtzu(X7, D7).
		Fcvtzu(w1, S7).
		Fmadd(D0, D1, D2, D3).
		Fmadd(S0, S1, S2, S3).
		Fnmsub(D4, D5, D6, D7).
		Fnmsub(S4, S5, S6, S7).
		Fcmp(D0, D1).
		Fcmp(S0, S1).
		FcmpZero(D2).
		FcmpZero(S2).
		LdrF(D3, X29, 16).
		LdrF(S3, X29, 8).
		StrF(D4, SP, 16).
		StrF(S4, SP, 8))

	b := arch.New()
	require.Equal(t, []uint32{
		enc(b.Fadd(D0, D1, D2)),
		enc(b.Fadd(S0, S1, S2)),
		enc(b.Fsub(D3, D4, D5)),
		enc(b.Fmul(S3, S4, S5)),
		enc(b.Fdiv(D6, D7, D0)),
		enc(b.Fmax(S6, S7, S0)),
		enc(b.Fmin(D1, D2, D3)),
		enc(b.Fneg(D4, D5)),
		enc(b.Fneg(S4, S5)),
		enc(b.Fmov(D6, D7)),
		enc(b.Fmov(S6, S7)),
		enc(b.FmovFromGpr(D0, X0)),
		enc(b.FmovFromGpr(S1, w1)),
		enc(b.FmovFromGpr(D2, XZR)),
		enc(b.FmovToGpr(X3, D5)),
		enc(b.FmovToGpr(w0, S7)),
		enc(b.FmovImm(D0, 0.5)),
		enc(b.FmovImm(S0, 1.5)),
		enc(b.Fcvt(S0, D0)),
		enc(b.Fcvt(D1, S1)),
		enc(b.Scvtf(D2, w1)),
		enc(b.Scvtf(D3, X2)),
		enc(b.Scvtf(S2, w1)),
		enc(b.Scvtf(S3, X3)),
		enc(b.Ucvtf(D4, X4)),
		enc(b.Ucvtf(S4, w0)),
		enc(b.Fcvtzs(X5, D5)),
		enc(b.Fcvtzs(w1, D6)),
		enc(b.Fcvtzs(X6, S5)),
		enc(b.Fcvtzs(w0, S6)),
		enc(b.Fcvtzu(X7, D7)),
		enc(b.Fcvtzu(w1, S7)),
		enc(b.Fmadd(D0, D1, D2, D3)),
		enc(b.Fmadd(S0, S1, S2, S3)),
		enc(b.Fnmsub(D4, D5, D6, D7)),
		enc(b.Fnmsub(S4, S5, S6, S7)),
		enc(b.Fcmp(D0, D1)),
		enc(b.Fcmp(S0, S1)),
		enc(b.FcmpZero(D2)),
		enc(b.FcmpZero(S2)),
		enc(b.LdrF(D3, X29, 16)),
		enc(b.LdrF(S3, X29, 8)),
		enc(b.StrF(D4, SP, 16)),
		enc(b.StrF(S4, SP, 8)),
	}, words(res.Code))
}

// TestTwinsFPClangPinned — the chain's FP words byte-compared with the
// constants the Apple clang (LLVM 17) oracle emits for the same source
// lines (the byte-parity pin of the whole FP surface).
func TestTwinsFPClangPinned(t *testing.T) {
	w2, w21, w23, w25 := wreg(t, 2), wreg(t, 21), wreg(t, 23), wreg(t, 25)
	x21, x23, x25 := reg(21), reg(23), reg(25)
	d8, d11 := dreg(8), dreg(11)
	d13 := dreg(13)
	d15, d16, d17 := dreg(15), dreg(16), dreg(17)
	d18, d19, d20 := dreg(18), dreg(19), dreg(20)
	d22, d24, d26 := dreg(22), dreg(24), dreg(26)
	s9, s12, s14 := sreg(9), sreg(12), sreg(14)
	s15, s16, s17 := sreg(15), sreg(16), sreg(17)
	s18, s19, s20 := sreg(18), sreg(19), sreg(20)
	s22, s24, s26 := sreg(22), sreg(24), sreg(26)
	d27, d28, d29, d30 := dreg(27), dreg(28), dreg(29), dreg(30)
	s27, s28, s29, s30 := sreg(27), sreg(28), sreg(29), sreg(30)

	res := assemble(t, New().
		Fadd(D0, D1, D2).
		Fadd(S0, S1, S2).
		Fsub(D0, D1, D2).
		Fsub(S0, S1, S2).
		Fmul(D0, D1, D2).
		Fmul(S0, S1, S2).
		Fdiv(D0, D1, D2).
		Fdiv(S0, S1, S2).
		Fmax(D0, D1, D2).
		Fmax(S0, S1, S2).
		Fmin(D0, D1, D2).
		Fmin(S0, S1, S2).
		Fneg(D3, D4).
		Fneg(S3, S4).
		Fmov(D5, D6).
		Fmov(S5, S6).
		FmovFromGpr(D7, X0).
		FmovToGpr(X1, d8).
		FmovFromGpr(s9, w2).
		FmovImm(d11, 0.5).
		FmovImm(s12, 0.5).
		FmovImm(d13, 1.5).
		FmovImm(s14, 1.5).
		Fcmp(d15, d16).
		Fcmp(s15, s16).
		FcmpZero(d17).
		FcmpZero(s17).
		Fcvt(s18, d19).
		Fcvt(d18, s19).
		Scvtf(d20, w21).
		Scvtf(d20, x21).
		Scvtf(s20, w21).
		Scvtf(s20, x21).
		Ucvtf(d22, w21).
		Ucvtf(d22, x21).
		Ucvtf(s22, w21).
		Ucvtf(s22, x21).
		Fcvtzs(w23, d24).
		Fcvtzs(x23, d24).
		Fcvtzs(w23, s24).
		Fcvtzs(x23, s24).
		Fcvtzu(w25, d26).
		Fcvtzu(x25, d26).
		Fcvtzu(w25, s26).
		Fcvtzu(x25, s26).
		Fmadd(d27, d28, d29, d30).
		Fmadd(s27, s28, s29, s30).
		Fnmsub(d27, d28, d29, d30).
		Fnmsub(s27, s28, s29, s30).
		LdrF(D0, X1, 8).
		LdrF(S0, X1, 4).
		StrF(D2, X1, 16).
		StrF(S2, X1, 8))

	require.Equal(t, []uint32{
		0x1e622820, 0x1e222820, // fadd d/s
		0x1e623820, 0x1e223820, // fsub d/s
		0x1e620820, 0x1e220820, // fmul d/s
		0x1e621820, 0x1e221820, // fdiv d/s
		0x1e624820, 0x1e224820, // fmax d/s
		0x1e625820, 0x1e225820, // fmin d/s
		0x1e614083, 0x1e214083, // fneg d/s
		0x1e6040c5, 0x1e2040c5, // fmov d/d, s/s
		0x9e670007, 0x9e660101, // fmov d<-x, x<-d
		0x1e270049,             // fmov s<-w
		0x1e6c100b, 0x1e2c100c, // fmov imm 0.5 d/s
		0x1e6f100d, 0x1e2f100e, // fmov imm 1.5 d/s
		0x1e7021e0, 0x1e3021e0, // fcmp d/d, s/s
		0x1e602228, 0x1e202228, // fcmp #0.0 d/s
		0x1e624272, 0x1e22c272, // fcvt s<-d, d<-s
		0x1e6202b4, 0x9e6202b4, 0x1e2202b4, 0x9e2202b4, // scvtf
		0x1e6302b6, 0x9e6302b6, 0x1e2302b6, 0x9e2302b6, // ucvtf
		0x1e780317, 0x9e780317, 0x1e380317, 0x9e380317, // fcvtzs
		0x1e790359, 0x9e790359, 0x1e390359, 0x9e390359, // fcvtzu
		0x1f5d7b9b, 0x1f1d7b9b, // fmadd d/s
		0x1f7dfb9b, 0x1f3dfb9b, // fnmsub d/s
		0xfd400420, 0xbd400420, // ldr d/s
		0xfd000822, 0xbd000822, // str d/s
	}, words(res.Code))
}

func TestTwinsFPErrors(t *testing.T) {
	// mixed kinds: the ctor error surfaces at Build
	_, buildErrs := New().Fadd(D0, S1, D2).Build()
	require.Len(t, buildErrs, 1)
	require.Contains(t, buildErrs[0].Error(), "fadd")

	// width-mismatched GPR move
	_, buildErrs = New().FmovFromGpr(D0, wreg(t, 0)).Build()
	require.Len(t, buildErrs, 1)
	require.Contains(t, buildErrs[0].Error(), "fmov")

	// fcvt needs differing kinds
	_, buildErrs = New().Fcvt(D0, D1).Build()
	require.Len(t, buildErrs, 1)
	require.Contains(t, buildErrs[0].Error(), "fcvt")

	// misaligned FP load offset
	_, buildErrs = New().LdrF(D0, X1, 3).Build()
	require.Len(t, buildErrs, 1)
	require.Contains(t, buildErrs[0].Error(), "ldr")

	// a non-encodable immediate builds fine and fails at Encode:
	// 0.1 is not one of the 256 VFP values
	bin, buildErrs := New().FmovImm(D0, 0.1).Build()
	require.Empty(t, buildErrs)
	errs := bin.Assemble(0).Errs
	require.Len(t, errs, 1)
	require.Contains(t, errs[0].Error(), "fmov")
}

// TestTwinsSimdClangPinned - the chain's SIMD words byte-compared with
// the constants the Apple clang (LLVM 17) oracle emits for the same
// source lines (the byte-parity pin of the decomposed SIMD families;
// the Builder is pinned transitively - every twin delegates to it).
func TestTwinsSimdClangPinned(t *testing.T) {
	w0 := wreg(t, 0)
	w1 := wreg(t, 1)

	res := assemble(t, New().
		And(V0, V1, V2, "16b").
		And(V3, V4, V5, "8b").
		Bic(V0, V1, V2, "16b").
		Orr(V0, V1, V2, "16b").
		Orr(V3, V4, V5, "8b").
		Orn(V0, V1, V2, "16b").
		Eor(V0, V1, V2, "16b").
		Bsl(V0, V1, V2, "16b").
		Bit(V0, V1, V2, "16b").
		Bif(V0, V1, V2, "16b").
		Add(V0, V1, V2, "16b").
		Add(V3, V4, V5, "8b").
		Cmeq(V0, V1, V2, "16b").
		Addp(V0, V1, V2, "16b").
		Sqrshl(V0, V1, V2, "16b").
		Cnt(V0, V1, "8b").
		Cnt(V0, V1, "16b").
		Rev32V(V0, V1, "8b").
		Not(V0, V1, "16b").
		Abs(V0, V1, "8b").
		RbitV(V0, V1, "16b").
		Shl(V0, V1, "16b", 3).
		Shl(V3, V4, "8b", 5).
		Sri(V0, V1, "16b", 4).
		Ushr(V0, V1, "16b", 5).
		Sshr(V0, V1, "16b", 5).
		Aese(V0, V1).
		Aesmc(V0, V1).
		Dup(V0, w0, "16b").
		DupElem(V0, V1, "4s", 2).
		Ins(V0, 1, w1, "s").
		Smov(X0, V0, "s", 1).
		Umov(w0, V0, "s", 1).
		InsElem(V0, V1, "d", 0, 0).
		InsElem(V2, V3, "s", 0, 0).
		MovSimd(V0, V1, "16b").
		MovSimd(V2, V3, "8b").
		Tbl(V0, V1, V2).
		Saddw(V0, V1, V2, "8h").
		Uaddw(V0, V1, V2, "4s").
		Usubw(V0, V1, V2, "8h").
		MlaElem(V0, V1, V2, "4s", 1).
		MlsElem(V0, V1, V2, "4s", 1).
		MulElem(V0, V1, V2, "4s", 1).
		SqdmulhElem(V0, V1, V2, "8h", 3).
		SqrdmulhElem(V0, V1, V2, "4s", 0).
		SqrdmlahElem(V0, V1, V2, "4s", 0).
		SqrdmlshElem(V0, V1, V2, "4s", 0).
		SmlalElem(V0, V1, V2, "4s", true, 1).
		FmlaElem(V0, V1, V2, "4s", 1).
		FmlsElem(V0, V1, V2, "2s", 0).
		FmulElem(V0, V1, V2, "2d", 0).
		FmulxElem(V0, V1, V2, "4s", 2).
		FcmlaElem(V0, V1, V2, "4s", 0, 90).
		InsElem(V0, V1, "s", 1, 2))

	require.Equal(t, []uint32{
		0x4e221c20,
		0x0e251c83,
		0x4e621c20,
		0x4ea21c20,
		0x0ea51c83,
		0x4ee21c20,
		0x6e221c20,
		0x6e621c20,
		0x6ea21c20,
		0x6ee21c20,
		0x4e228420,
		0x0e258483,
		0x6e228c20,
		0x4e22bc20,
		0x4e225c20,
		0x0e205820,
		0x4e205820,
		0x2e200820,
		0x6e205820,
		0x0e20b820,
		0x6e605820,
		0x4f0b5420,
		0x0f0d5483,
		0x6f0c4420,
		0x6f0b0420,
		0x4f0b0420,
		0x4e284820,
		0x4e286820,
		0x4e010c00,
		0x4e140420,
		0x4e0c1c20,
		0x4e0c2c00,
		0x0e0c3c00,
		0x6e080420,
		0x6e040462,
		0x4ea11c20,
		0x0ea31c62,
		0x4e020020,
		0x4e221020,
		0x6e621020,
		0x6e223020,
		0x6fa20020,
		0x6fa24020,
		0x4fa28020,
		0x4f72c020,
		0x4f82d020,
		0x6f82d020,
		0x6f82f020,
		0x4f522020,
		0x4fa21020,
		0x0f825020,
		0x4fc29020,
		0x6f829820,
		0x6f823020,
		0x6e0c4420,
	}, words(res.Code))
}

func TestTwinsSimdErrors(t *testing.T) {
	// the logical group takes only .8b/.16b (bits 23:22 are its opcode)
	_, buildErrs := New().And(V0, V1, V2, "4h").Build()
	require.Len(t, buildErrs, 1)
	require.Contains(t, buildErrs[0].Error(), "and")

	// a lane index beyond the lane count
	_, buildErrs = New().DupElem(V0, V1, "4s", 4).Build()
	require.Len(t, buildErrs, 1)
	require.Contains(t, buildErrs[0].Error(), "dup")

	// umov: the element must fill the destination register
	_, buildErrs = New().Umov(X0, V0, "s", 1).Build()
	require.Len(t, buildErrs, 1)
	require.Contains(t, buildErrs[0].Error(), "umov")

	// a shift out of the lane width's range
	_, buildErrs = New().Shl(V0, V1, "16b", 16).Build()
	require.Len(t, buildErrs, 1)
	require.Contains(t, buildErrs[0].Error(), "shl")

	// an fp arrangement on an integer by-element form
	_, buildErrs = New().MlaElem(V0, V1, V2, "2d", 1).Build()
	require.Len(t, buildErrs, 1)
	require.Contains(t, buildErrs[0].Error(), "mla")

	// fcmla rotation is one of #0/#90/#180/#270
	_, buildErrs = New().FcmlaElem(V0, V1, V2, "4s", 0, 45).Build()
	require.Len(t, buildErrs, 1)
	require.Contains(t, buildErrs[0].Error(), "fcmla")
}
