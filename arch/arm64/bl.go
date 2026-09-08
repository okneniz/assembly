package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Bl — bl off (call, imm26).
type Bl struct {
	base

	off imm // pc-relative byte offset
}

const blMatch = 0x94000000

// Bl — bl off: off — the pc-relative byte offset of the call
// destination (the ±128MB imm26 range is checked at encode time; the
// absolute target is off + the instruction address).
func (Builder) Bl(off int64) Instr {
	return Bl{off: immNum(off)}
}

func decodeBl(w uint32) Instr {
	return Bl{
		base: newBase(w),
		off:  immNum(signExtendN(w&0x3ffffff, 26) * 4),
	}
}

func (i Bl) ObjDump(ctx disasm.ViewCtx) string {
	target := immNum(int64(ctx.Addr()) + i.off.val)
	return "bl " + target.textHex()
}

func (i Bl) Encode(w io.Writer) (int64, error) {
	bits, err := offBits(i.off.val, 26)
	if err != nil {
		return 0, fmt.Errorf("bl: %w", err)
	}

	return writeWord(w, blMatch|bits)
}
