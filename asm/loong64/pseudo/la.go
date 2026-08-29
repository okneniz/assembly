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

// La is la rd, sym: a pcalau12i+addi.d pair (a fixed 8 bytes). An
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

// resolveLa is the evaluator wired to parsing: la rd, sym.
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
