package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Bcond — b.cond off (imm19; cond — in the base word).
type Bcond struct {
	base

	cond string
	off  imm // pc-relative byte offset
}

func decodeBcondOf(cond string) func(uint32) Instr {
	return func(w uint32) Instr {
		return Bcond{
			base: newBase(w),
			cond: cond,
			off:  immNum(signExtendN(w>>5&0x7ffff, 19) * 4),
		}
	}
}

func (i Bcond) ObjDump(ctx disasm.ViewCtx) string {
	target := immNum(int64(ctx.Addr()) + i.off.val)
	return fmt.Sprintf("b.%s %s", i.cond, target.textHex())
}

func (i Bcond) Encode(w io.Writer) (int64, error) {
	c, err := condNum(i.cond)
	if err != nil {
		return 0, fmt.Errorf("b.cond: %w", err)
	}

	bits, err := offBits(i.off.val, 19)
	if err != nil {
		return 0, fmt.Errorf("b.%s: %w", i.cond, err)
	}

	return writeWord(w, 0x54000000|c|bits<<5)
}

// Bcond — b.cond off: off — the pc-relative byte offset of the branch
// destination (the ±1MB imm19 range is checked at encode time; the
// absolute target is off + the instruction address); cond — the
// standard condition names (eq/ne/.../nv, see condNum).
func (Builder) Bcond(cond string, off int64) (Instr, error) {
	if _, err := condNum(cond); err != nil {
		return nil, fmt.Errorf("arm64.NewBcond: operand cond: %w", err)
	}

	return Bcond{cond: cond, off: immNum(off)}, nil
}
