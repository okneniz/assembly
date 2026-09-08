package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// B — b off (unconditional branch, imm26, ±128MB).
type B struct {
	base

	off imm // pc-relative byte offset
}

const bMatch = 0x14000000

// B — b off (the pc-relative byte offset; the absolute target is off +
// the instruction address).
func (Builder) B(off int64) Instr {
	return B{
		off: immNum(off),
	}
}

func decodeB(w uint32) Instr {
	return B{
		base: newBase(w),
		off:  immNum(signExtendN(w&0x3ffffff, 26) * 4),
	}
}

func (i B) ObjDump(ctx disasm.ViewCtx) string {
	target := immNum(int64(ctx.Addr()) + i.off.val)
	return "b " + target.textHex()
}

func (i B) Encode(w io.Writer) (int64, error) {
	bits, err := offBits(i.off.val, 26)
	if err != nil {
		return 0, fmt.Errorf("b: %w", err)
	}

	return writeWord(w, bMatch|bits)
}
