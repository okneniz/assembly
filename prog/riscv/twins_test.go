package riscv

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"

	arch "github.com/okneniz/assembly/arch/riscv"
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

// enc - the word of a directly built instruction (the NoRVC form the
// chain programs encode in): the oracle the twins are byte-compared
// against. Panics on an encode error - the oracle operands are valid
// by construction (the pre-minted-reg pattern of regs.go).
func enc(i arch.Instr) uint32 {
	var buf bytes.Buffer
	if _, err := i.Encode(&buf, arch.EncOpts{NoRVC: true}); err != nil {
		panic(err)
	}

	return binary.LittleEndian.Uint32(buf.Bytes())
}

// offv - a validated load/store offset for the oracle side. Panics on
// a range error - the values below are valid by construction.
func offv(v int64) arch.Off {
	o, err := arch.New().Off(v)
	if err != nil {
		panic(err)
	}

	return o
}

func TestTwinsArithmetic(t *testing.T) {
	res := assemble(t, New().
		Add(T0, T1, T2).
		Addw(A0, A1, A2).
		Addi(T0, T1, -5).
		Addiw(A1, A2, 100).
		Sub(T2, A0, A1).
		Subw(A2, T0, T1).
		Mul(T0, T1, T2).
		Mulh(A0, A1, A2).
		Mulhsu(T1, T2, A0).
		Mulhu(A2, T0, T1).
		Mulw(T2, A0, A1).
		Div(T0, T1, T2).
		Divu(A0, A1, A2).
		Divw(T1, T2, A0).
		Divuw(A2, T0, T1).
		Rem(T2, A0, A1).
		Remu(T0, T1, T2).
		Remw(A0, A1, A2).
		Remuw(T1, T2, A0).
		Slt(A2, T0, T1).
		Sltu(T2, A0, A1).
		Slti(T0, T1, -3).
		Sltiu(A1, A2, 7))

	b := arch.New()
	v12a, err := b.Imm12(-5)
	require.NoError(t, err)
	v12b, err := b.Imm12(100)
	require.NoError(t, err)
	v12c, err := b.Imm12(-3)
	require.NoError(t, err)
	v12d, err := b.Imm12(7)
	require.NoError(t, err)

	require.Equal(t, []uint32{
		enc(b.Add(T0, T1, T2)),
		enc(b.Addw(A0, A1, A2)),
		enc(b.Addi(T0, T1, v12a)),
		enc(b.Addiw(A1, A2, v12b)),
		enc(b.Sub(T2, A0, A1)),
		enc(b.Subw(A2, T0, T1)),
		enc(b.Mul(T0, T1, T2)),
		enc(b.Mulh(A0, A1, A2)),
		enc(b.Mulhsu(T1, T2, A0)),
		enc(b.Mulhu(A2, T0, T1)),
		enc(b.Mulw(T2, A0, A1)),
		enc(b.Div(T0, T1, T2)),
		enc(b.Divu(A0, A1, A2)),
		enc(b.Divw(T1, T2, A0)),
		enc(b.Divuw(A2, T0, T1)),
		enc(b.Rem(T2, A0, A1)),
		enc(b.Remu(T0, T1, T2)),
		enc(b.Remw(A0, A1, A2)),
		enc(b.Remuw(T1, T2, A0)),
		enc(b.Slt(A2, T0, T1)),
		enc(b.Sltu(T2, A0, A1)),
		enc(b.Slti(T0, T1, v12c)),
		enc(b.Sltiu(A1, A2, v12d)),
	}, words(res.Code))
}

func TestTwinsLogicalShifts(t *testing.T) {
	res := assemble(t, New().
		And(T0, T1, T2).
		Andi(A0, A1, -2).
		Or(T1, T2, A0).
		Ori(A2, T0, 5).
		Xor(T2, A0, A1).
		Xori(T0, T1, 0xff).
		Sll(A0, A1, A2).
		Sllw(T1, T2, A0).
		Slli(A2, T0, 4).
		Slliw(T2, A0, 5).
		Srl(T0, T1, T2).
		Srlw(A0, A1, A2).
		Srli(T1, T2, 6).
		Srliw(A2, T0, 7).
		Sra(T0, T1, T2).
		Sraw(A0, A1, A2).
		Srai(T1, T2, 8).
		Sraiw(A2, T0, 9))

	b := arch.New()
	vM2, err := b.Imm12(-2)
	require.NoError(t, err)
	v5, err := b.Imm12(5)
	require.NoError(t, err)
	vFF, err := b.Imm12(0xff)
	require.NoError(t, err)
	v4, err := b.Imm12(4)
	require.NoError(t, err)
	v5s, err := b.Imm12(5)
	require.NoError(t, err)
	v6, err := b.Imm12(6)
	require.NoError(t, err)
	v7, err := b.Imm12(7)
	require.NoError(t, err)
	v8, err := b.Imm12(8)
	require.NoError(t, err)
	v9, err := b.Imm12(9)
	require.NoError(t, err)

	require.Equal(t, []uint32{
		enc(b.And(T0, T1, T2)),
		enc(b.Andi(A0, A1, vM2)),
		enc(b.Or(T1, T2, A0)),
		enc(b.Ori(A2, T0, v5)),
		enc(b.Xor(T2, A0, A1)),
		enc(b.Xori(T0, T1, vFF)),
		enc(b.Sll(A0, A1, A2)),
		enc(b.Sllw(T1, T2, A0)),
		enc(b.Slli(A2, T0, v4)),
		enc(b.Slliw(T2, A0, v5s)),
		enc(b.Srl(T0, T1, T2)),
		enc(b.Srlw(A0, A1, A2)),
		enc(b.Srli(T1, T2, v6)),
		enc(b.Srliw(A2, T0, v7)),
		enc(b.Sra(T0, T1, T2)),
		enc(b.Sraw(A0, A1, A2)),
		enc(b.Srai(T1, T2, v8)),
		enc(b.Sraiw(A2, T0, v9)),
	}, words(res.Code))
}

func TestTwinsUpperJumps(t *testing.T) {
	res := assemble(t, New().
		Lui(T0, 0x12345).
		Auipc(A0, 0x54321>>12&0xfffff).
		Jalr(Ra, T0, 0).
		Jalr(T2, A0, 8).
		Fence(3))

	b := arch.New()
	v20a, err := b.Imm20(0x12345)
	require.NoError(t, err)
	v20b, err := b.Imm20(0x54321 >> 12 & 0xfffff)
	require.NoError(t, err)

	require.Equal(t, []uint32{
		enc(b.Lui(T0, v20a)),
		enc(b.Auipc(A0, v20b)),
		enc(b.Jalr(Ra, T0, offv(0))),
		enc(b.Jalr(T2, A0, offv(8))),
		enc(b.Fence(3)),
	}, words(res.Code))
}

func TestTwinsLoadsStores(t *testing.T) {
	fa0 := reg(10)
	fa1 := reg(11)
	fa2 := reg(12)

	res := assemble(t, New().
		Lb(T0, A0, 8).
		Lbu(T1, A0, 9).
		Lh(T2, A0, 10).
		Lhu(A1, A0, 12).
		Lw(A2, A0, 16).
		Lwu(T0, A0, 20).
		Ld(T1, A0, 24).
		Flw(fa0, A0, 28).
		Fld(fa1, A0, 32).
		Sb(T2, A0, 8).
		Sh(A1, A0, 10).
		Sw(A2, A0, 16).
		Sd(T0, A0, 24).
		Fsw(fa2, A0, 28).
		Fsd(fa0, A0, 32))

	b := arch.New()
	require.Equal(t, []uint32{
		enc(b.Lb(T0, A0, offv(8))),
		enc(b.Lbu(T1, A0, offv(9))),
		enc(b.Lh(T2, A0, offv(10))),
		enc(b.Lhu(A1, A0, offv(12))),
		enc(b.Lw(A2, A0, offv(16))),
		enc(b.Lwu(T0, A0, offv(20))),
		enc(b.Ld(T1, A0, offv(24))),
		enc(b.Flw(fa0, A0, offv(28))),
		enc(b.Fld(fa1, A0, offv(32))),
		enc(b.Sb(T2, A0, offv(8))),
		enc(b.Sh(A1, A0, offv(10))),
		enc(b.Sw(A2, A0, offv(16))),
		enc(b.Sd(T0, A0, offv(24))),
		enc(b.Fsw(fa2, A0, offv(28))),
		enc(b.Fsd(fa0, A0, offv(32))),
	}, words(res.Code))
}

func TestTwinsAmoCsr(t *testing.T) {
	res := assemble(t, New().
		AmoswapW(T0, A0, T1).
		AmoswapD(T1, A0, T2).
		AmoaddW(T2, A0, T0).
		AmoaddD(A0, A0, T1).
		AmoxorW(A1, A0, T2).
		AmoxorD(A2, A0, T0).
		AmoandW(T0, A1, T1).
		AmoandD(T1, A1, T2).
		AmoorW(T2, A1, T0).
		AmoorD(A0, A1, T1).
		AmominW(A1, A1, T2).
		AmominD(A2, A1, T0).
		AmominuW(T0, A2, T1).
		AmominuD(T1, A2, T2).
		AmomaxW(T2, A2, T0).
		AmomaxD(A0, A2, T1).
		AmomaxuW(A1, A2, T2).
		AmomaxuD(A2, A2, T0).
		Csrrw(T0, 0x341, T1).
		Csrrs(T1, 0x342, T2).
		Csrrc(T2, 0x300, A0).
		Csrrwi(A0, 0x105, 5).
		Csrrsi(A1, 0x106, 6).
		Csrrci(A2, 0x107, 7))

	b := arch.New()
	require.Equal(t, []uint32{
		enc(b.AmoswapW(T0, A0, T1)),
		enc(b.AmoswapD(T1, A0, T2)),
		enc(b.AmoaddW(T2, A0, T0)),
		enc(b.AmoaddD(A0, A0, T1)),
		enc(b.AmoxorW(A1, A0, T2)),
		enc(b.AmoxorD(A2, A0, T0)),
		enc(b.AmoandW(T0, A1, T1)),
		enc(b.AmoandD(T1, A1, T2)),
		enc(b.AmoorW(T2, A1, T0)),
		enc(b.AmoorD(A0, A1, T1)),
		enc(b.AmominW(A1, A1, T2)),
		enc(b.AmominD(A2, A1, T0)),
		enc(b.AmominuW(T0, A2, T1)),
		enc(b.AmominuD(T1, A2, T2)),
		enc(b.AmomaxW(T2, A2, T0)),
		enc(b.AmomaxD(A0, A2, T1)),
		enc(b.AmomaxuW(A1, A2, T2)),
		enc(b.AmomaxuD(A2, A2, T0)),
		enc(b.Csrrw(T0, 0x341, T1)),
		enc(b.Csrrs(T1, 0x342, T2)),
		enc(b.Csrrc(T2, 0x300, A0)),
		enc(b.Csrrwi(A0, 0x105, 5)),
		enc(b.Csrrsi(A1, 0x106, 6)),
		enc(b.Csrrci(A2, 0x107, 7)),
	}, words(res.Code))
}

func TestTwinsFP(t *testing.T) {
	fa0 := reg(10)
	fa1 := reg(11)
	fa2 := reg(12)
	fa3 := reg(13)

	res := assemble(t, New().
		FaddS(fa0, fa1, fa2, 0).
		FaddD(fa0, fa1, fa2, 0).
		FsubS(fa1, fa2, fa3, 1).
		FsubD(fa1, fa2, fa3, 1).
		FmulS(fa2, fa3, fa0, 2).
		FmulD(fa2, fa3, fa0, 2).
		FdivS(fa3, fa0, fa1, 3).
		FdivD(fa3, fa0, fa1, 3).
		FmaddS(fa0, fa1, fa2, fa3, 0).
		FmaddD(fa0, fa1, fa2, fa3, 0).
		FmsubS(fa1, fa2, fa3, fa0, 4).
		FmsubD(fa1, fa2, fa3, fa0, 4).
		FnmaddS(fa2, fa3, fa0, fa1, 5).
		FnmaddD(fa2, fa3, fa0, fa1, 5).
		FnmsubS(fa3, fa0, fa1, fa2, 6).
		FnmsubD(fa3, fa0, fa1, fa2, 6))

	b := arch.New()
	require.Equal(t, []uint32{
		enc(b.FaddS(fa0, fa1, fa2, 0)),
		enc(b.FaddD(fa0, fa1, fa2, 0)),
		enc(b.FsubS(fa1, fa2, fa3, 1)),
		enc(b.FsubD(fa1, fa2, fa3, 1)),
		enc(b.FmulS(fa2, fa3, fa0, 2)),
		enc(b.FmulD(fa2, fa3, fa0, 2)),
		enc(b.FdivS(fa3, fa0, fa1, 3)),
		enc(b.FdivD(fa3, fa0, fa1, 3)),
		enc(b.FmaddS(fa0, fa1, fa2, fa3, 0)),
		enc(b.FmaddD(fa0, fa1, fa2, fa3, 0)),
		enc(b.FmsubS(fa1, fa2, fa3, fa0, 4)),
		enc(b.FmsubD(fa1, fa2, fa3, fa0, 4)),
		enc(b.FnmaddS(fa2, fa3, fa0, fa1, 5)),
		enc(b.FnmaddD(fa2, fa3, fa0, fa1, 5)),
		enc(b.FnmsubS(fa3, fa0, fa1, fa2, 6)),
		enc(b.FnmsubD(fa3, fa0, fa1, fa2, 6)),
	}, words(res.Code))
}

func TestTwinsBranches(t *testing.T) {
	// forward and backward targets through the label machinery
	p := New().
		Label("top").        // @0x00
		Bge(T0, T1, "fwd").  // @0x00 → fwd @0x14: +20
		Bltu(T1, T2, "top"). // @0x04 → top @0x00: -4
		Blt(A0, A1, "fwd").  // @0x08 → fwd @0x14: +12
		Addi(Zero, Zero, 0). // @0x0c - padding (the nop encoding)
		Addi(Zero, Zero, 0). // @0x10
		Label("fwd").        // @0x14
		Entry("top")

	res := assemble(t, p)
	b := arch.New()
	v0, err := b.Imm12(0)
	require.NoError(t, err)

	require.Equal(t, []uint32{
		enc(b.Bge(T0, T1, 20)),
		enc(b.Bltu(T1, T2, -4)),
		enc(b.Blt(A0, A1, 12)),
		enc(b.Addi(Zero, Zero, v0)),
		enc(b.Addi(Zero, Zero, v0)),
	}, words(res.Code))
}

func TestTwinsErrors(t *testing.T) {
	// imm out of range: the deferred error surfaces at Build
	_, buildErrs := New().Addi(T0, T1, 4096).Build()
	require.Len(t, buildErrs, 1)
	require.Contains(t, buildErrs[0].Error(), "addi")

	// off out of range: the deferred error surfaces at Build
	_, buildErrs = New().Ld(T0, A0, 4096).Build()
	require.Len(t, buildErrs, 1)
	require.Contains(t, buildErrs[0].Error(), "ld")

	// undefined branch target surfaces at Assemble
	bin, buildErrs := New().Bge(T0, T1, "nowhere").Build()
	require.Empty(t, buildErrs)
	errs := bin.Assemble(0).Errs
	require.Len(t, errs, 1)
	require.Contains(t, errs[0].Error(), "nowhere")
}
