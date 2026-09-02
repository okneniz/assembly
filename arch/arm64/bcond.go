package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Bcond — b.cond target (imm19; cond — in the base word).
type Bcond struct {
	base

	cond   string
	target imm
}

// Bcond — b.cond target: target — the absolute address of the
// branch destination (the ±1MB imm19 range is checked at encode time,
// from pc); cond — the standard condition names (eq/ne/.../nv, see
// condNum).
func (Builder) Bcond(cond string, target int64) (Instr, error) {
	if _, err := condNum(cond); err != nil {
		return nil, fmt.Errorf("arm64.NewBcond: operand cond: %w", err)
	}

	return Bcond{cond: cond, target: immNum(target)}, nil
}

func decodeBcondOf(cond string) func(uint32, uint64) Instr {
	return func(w uint32, addr uint64) Instr {
		return Bcond{
			base:   newBase(addr, w),
			cond:   cond,
			target: immNum(int64(addr) + signExtendN(w>>5&0x7ffff, 19)*4),
		}
	}
}

func (i Bcond) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("b.%s %s", i.cond, i.target.textHex())
}

func (i Bcond) Encode(w io.Writer, pc uint64) (int64, error) {
	target := i.target.val

	c, err := condNum(i.cond)
	if err != nil {
		return 0, fmt.Errorf("b.cond: %w", err)
	}

	bits, err := brBits(target, int64(pc), 19)
	if err != nil {
		return 0, fmt.Errorf("b.%s: %w", i.cond, err)
	}

	return writeWord(w, 0x54000000|c|bits<<5)
}

func (i Bcond) MarshalJSON() ([]byte, error) {
	return i.marshal(
		"b."+i.cond,
		i.ObjDump(disasm.DefaultViewCtx()),
		"Branch",
		map[string]any{"cond": i.cond},
	)
}
