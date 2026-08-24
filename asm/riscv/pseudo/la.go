package pseudo

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	arch "github.com/okneniz/assembly/arch/riscv"
	"github.com/okneniz/assembly/asm"
	riscv "github.com/okneniz/assembly/asm/riscv"
	"github.com/okneniz/assembly/disasm"
)

// La is la rd, sym: an auipc+addi pair via pcrel (a fixed 8 bytes).
// An evaluated form: the target and address are known, the encoding
// is pure.
type La struct {
	rd     string
	target int64
	pc     uint64
}

func (i La) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("la %s, %#x", i.rd, i.target)
}

func (i La) Len() int {
	return 8
}
func (i La) Addr() uint64 {
	return 0
} // assembly side: there is no decode address

// Encode writes the pcrel auipc+addi pair without compression: the
// length is fixed in both passes regardless of symbol values (the
// placeholder yields rel=0).
func (i La) Encode(w io.Writer) (int64, error) {
	hi, lo := arch.PcrelHiLo(i.target - int64(i.pc))
	hiBits, err := arch.EncU(hi & 0xfffff)
	if err != nil {
		return 0, fmt.Errorf("la: %w", err)
	}

	loBits, err := arch.EncI(lo)
	if err != nil {
		return 0, fmt.Errorf("la: %w", err)
	}

	rd := arch.RegBits(i.rd)
	auipc := arch.EncodingWord("auipc") | rd<<7 | hiBits
	addi := arch.EncodingWord("addi") | rd<<7 | rd<<15 | loBits
	var buf bytes.Buffer
	if _, err := arch.WriteWord(&buf, auipc); err != nil {
		return 0, err
	}

	if _, err := arch.WriteWord(&buf, addi); err != nil {
		return 0, err
	}

	n, err := w.Write(buf.Bytes())
	return int64(n), err
}

func (i La) MarshalJSON() ([]byte, error) {
	return arch.MarshalDTO(
		arch.Base{},
		"la",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Pseudo",
		map[string]any{"rd": i.rd},
	)
}

// resolveLa is the evaluator wired to parsing: la rd, sym.
func resolveLa(ops []riscv.Op, ctx asm.Ctx) (asm.Resolved, error) {
	if len(ops) != 2 {
		return nil, errors.New("la: want rd, sym")
	}

	rd, err := riscv.WantReg(ops[0], false)
	if err != nil {
		return nil, fmt.Errorf("la: %w", err)
	}

	e, err := riscv.WantExpr(ops[1])
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
