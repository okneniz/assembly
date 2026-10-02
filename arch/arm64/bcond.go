package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Bcond — b.cond off (imm19; cond — in the base word).
type Bcond struct {
	cond string
	off  imm // pc-relative byte offset
}

// newBcond - the Bcond constructor: the struct is assembled only here.
func newBcond(cond string, off imm) Bcond {
	return Bcond{
		cond: cond,
		off:  off,
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
