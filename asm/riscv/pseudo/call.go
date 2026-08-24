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

// Call is call sym: auipc ra, hi + jalr ra, lo(ra) (a fixed 8 bytes).
// An evaluated form: the target and address are known, the encoding
// is pure.
type Call struct {
	target int64
	pc     uint64
}

func (i Call) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("call %#x", i.target)
}

func (i Call) Len() int {
	return 8
}
func (i Call) Addr() uint64 {
	return 0
} // assembly side: there is no decode address

func (i Call) Encode(w io.Writer) (int64, error) {
	hi, lo := arch.PcrelHiLo(i.target - int64(i.pc))
	hiBits, err := arch.EncU(hi & 0xfffff)
	if err != nil {
		return 0, fmt.Errorf("call: %w", err)
	}

	loBits, err := arch.EncI(lo)
	if err != nil {
		return 0, fmt.Errorf("call: %w", err)
	}

	auipc := arch.EncodingWord("auipc") | 1<<7 | hiBits // ra
	jalr := arch.EncodingWord("jalr") | 1<<7 | 1<<15 | loBits
	var buf bytes.Buffer
	if _, err := arch.WriteWord(&buf, auipc); err != nil {
		return 0, err
	}

	if _, err := arch.WriteWord(&buf, jalr); err != nil {
		return 0, err
	}

	n, err := w.Write(buf.Bytes())
	return int64(n), err
}

func (i Call) MarshalJSON() ([]byte, error) {
	return arch.MarshalDTO(arch.Base{}, "call", i.ObjDump(disasm.DefaultViewCtx()), "Pseudo", nil)
}

// resolveCall is the evaluator wired to parsing: call sym.
func resolveCall(ops []riscv.Op, ctx asm.Ctx) (asm.Resolved, error) {
	if len(ops) != 1 {
		return nil, errors.New("call: want sym")
	}

	e, err := riscv.WantExpr(ops[0])
	if err != nil {
		return nil, fmt.Errorf("call: %w", err)
	}

	t, terr := e.Eval(ctx.Resolve)
	if terr != nil {
		return nil, fmt.Errorf("call: %w", terr)
	}

	return Call{
		target: t,
		pc:     ctx.Addr(),
	}, nil
}
