package pseudo

// Tail is tail sym: jal zero, sym (not compressed: RVC c.j only for
// close numeric targets within +/-2KB; a tail call reaches farther).
// An evaluated form.

import (
	"errors"
	"fmt"
	"io"

	arch "github.com/okneniz/assembly/arch/riscv"
	"github.com/okneniz/assembly/asm"
	riscv "github.com/okneniz/assembly/asm/riscv"
	"github.com/okneniz/assembly/disasm"
)

type Tail struct {
	target int64
	pc     uint64
}

func (i Tail) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("tail %#x", i.target)
}

func (i Tail) Len() int {
	return 4
}
func (i Tail) Addr() uint64 {
	return 0
} // assembly side: there is no decode address

func (i Tail) Encode(w io.Writer) (int64, error) {
	bits, err := arch.EncJ(i.target - int64(i.pc))
	if err != nil {
		return 0, fmt.Errorf("tail: %w", err)
	}

	return arch.WriteWord(w, arch.EncodingWord("jal")|bits)
}

func (i Tail) MarshalJSON() ([]byte, error) {
	return arch.MarshalDTO(arch.Base{}, "tail", i.ObjDump(disasm.DefaultViewCtx()), "Pseudo", nil)
}

// resolveTail is the evaluator wired to parsing: tail sym.
func resolveTail(ops []riscv.Op, ctx asm.Ctx) (asm.Resolved, error) {
	if len(ops) != 1 {
		return nil, errors.New("tail: want sym")
	}

	e, err := riscv.WantExpr(ops[0])
	if err != nil {
		return nil, fmt.Errorf("tail: %w", err)
	}

	t, terr := e.Eval(ctx.Resolve)
	if terr != nil {
		return nil, fmt.Errorf("tail: %w", terr)
	}

	return Tail{
		target: t,
		pc:     ctx.Addr(),
	}, nil
}
