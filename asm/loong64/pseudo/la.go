package pseudo

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	arch "github.com/okneniz/assembly/arch/loong64"
	"github.com/okneniz/assembly/asm"
	loong64 "github.com/okneniz/assembly/asm/loong64"
)

// La is la rd, sym: a pcalau12i+addi.d pair (a fixed 8 bytes) - the
// la.pcrel form in llvm-mc terms (its plain la is the GOT-indirect
// pcalau12i+ld.d pair, which needs a linker-built GOT entry; this
// assembler resolves everything absolutely, so la IS la.pcrel here). An
// evaluated form: the target and address are known, the encoding is
// pure. pcalau12i clears the low 12 bits AFTER the page addition, so
// the split is computed against the page-aligned pc:
//
//	D  = target - (pc &^ 0xfff)
//	lo = sext12(D)        (the addi.d immediate)
//	hi = (D - lo) >> 12   (the pcalau12i immediate)
type La struct {
	rd     string
	target int64
	pc     uint64
}

// Encode writes the pcalau12i+addi.d pair.
func (i La) Encode(w io.Writer) (int64, error) {
	page := int64(i.pc &^ 0xfff)
	d := i.target - page
	lo := d & 0xfff
	if lo >= 0x800 {
		lo -= 0x1000
	}

	hi := (d - lo) >> 12
	if hi < -(1<<19) || hi > (1<<19)-1 {
		return 0, fmt.Errorf(
			"la: offset %#x does not fit the pcalau12i+addi.d pair",
			i.target-int64(i.pc),
		)
	}

	rd, err := arch.RegNumOf(i.rd)
	if err != nil {
		return 0, fmt.Errorf("la: %w", err)
	}

	pcalau := arch.EncodingWord("pcalau12i") | rd | uint32(hi&(1<<20-1))<<5
	addi := arch.EncodingWord("addi.d") | rd | rd<<5 | uint32(lo&(1<<12-1))<<10

	var buf bytes.Buffer
	if _, err := arch.WriteWord(&buf, pcalau); err != nil {
		return 0, err
	}

	if _, err := arch.WriteWord(&buf, addi); err != nil {
		return 0, err
	}

	n, err := w.Write(buf.Bytes())
	return int64(n), err
}

// resolveLa is the evaluator wired to parsing: la/la.pcrel rd, sym.
func resolveLa(ops []loong64.Op, ctx asm.Ctx) (asm.Resolved, error) {
	if len(ops) != 2 {
		return nil, errors.New("la: want rd, sym")
	}

	rd, err := loong64.WantReg(ops[0])
	if err != nil {
		return nil, fmt.Errorf("la: %w", err)
	}

	e, err := loong64.WantExpr(ops[1])
	if err != nil {
		return nil, fmt.Errorf("la: %w", err)
	}

	t, terr := e.Eval(ctx.Resolve)
	if terr != nil {
		return nil, fmt.Errorf("la: %w", terr)
	}

	return La{
		rd:     rd,
		target: t,
		pc:     ctx.Addr(),
	}, nil
}

// LaAbs is la.abs rd, sym: the 64-bit absolute ladder lu12i.w+ori+
// lu32i.d+lu52i.d (a fixed 16 bytes) - the same four words llvm-mc
// expands la.abs into (and lld keeps them; no relaxation):
//
//	lu12i.w rd, (v + 0x800) >> 12        bits [31:12] (ori carry)
//	ori     rd, rd, v & 0xfff            bits [11:0]
//	lu32i.d rd, bits[51:32]              (raw 20-bit field)
//	lu52i.d rd, rd, bits[63:52]          (raw 12-bit field)
type LaAbs struct {
	rd     string
	target int64
}

// Encode writes the lu12i.w+ori+lu32i.d+lu52i.d ladder.
func (i LaAbs) Encode(w io.Writer) (int64, error) {
	v := i.target
	hi := (v + 0x800) >> 12
	lo := v & 0xfff
	w32 := (v >> 32) & (1<<20 - 1)
	w52 := (v >> 52) & 0xfff

	rd, err := arch.RegNumOf(i.rd)
	if err != nil {
		return 0, fmt.Errorf("la.abs: %w", err)
	}

	words := []uint32{
		arch.EncodingWord("lu12i.w") | rd | uint32(hi&(1<<20-1))<<5,
		arch.EncodingWord("ori") | rd | rd<<5 | uint32(lo)<<10,
		arch.EncodingWord("lu32i.d") | rd | uint32(w32)<<5,
		arch.EncodingWord("lu52i.d") | rd | rd<<5 | uint32(w52)<<10,
	}

	var buf bytes.Buffer
	for _, word := range words {
		if _, err := arch.WriteWord(&buf, word); err != nil {
			return 0, err
		}
	}

	n, err := w.Write(buf.Bytes())
	return int64(n), err
}

// resolveLaAbs is the evaluator wired to parsing: la.abs rd, sym.
func resolveLaAbs(ops []loong64.Op, ctx asm.Ctx) (asm.Resolved, error) {
	if len(ops) != 2 {
		return nil, errors.New("la.abs: want rd, sym")
	}

	rd, err := loong64.WantReg(ops[0])
	if err != nil {
		return nil, fmt.Errorf("la.abs: %w", err)
	}

	e, err := loong64.WantExpr(ops[1])
	if err != nil {
		return nil, fmt.Errorf("la.abs: %w", err)
	}

	t, terr := e.Eval(ctx.Resolve)
	if terr != nil {
		return nil, fmt.Errorf("la.abs: %w", terr)
	}

	return LaAbs{rd: rd, target: t}, nil
}

// resolveLaGot rejects la.got explicitly: the GOT-indirect form
// (pcalau12i+ld.d through a linker-created GOT slot) cannot resolve in
// this standalone assembler (no relocations, everything absolute).
func resolveLaGot(_ []loong64.Op, _ asm.Ctx) (asm.Resolved, error) {
	return nil, errors.New("la.got: the GOT form needs a linker-built GOT slot" +
		" (this assembler resolves absolutely; use la/la.pcrel or la.abs)")
}
