package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// B — b off (unconditional branch, imm26, ±128MB).
type B struct {
	off imm // pc-relative byte offset
}

// newB - the B constructor: the struct is assembled only
// here (the Builder method and the decoder call it).
func newB(off imm) (B, error) { //nolint:unparam // uniform (Instr, error) decodeCtor type
	return B{
		off: off,
	}, nil
}

const bMatch = 0x14000000

func (i B) Encode(w io.Writer) (int64, error) {
	bits, err := offBits(i.off.val, 26)
	if err != nil {
		return 0, fmt.Errorf("b: %w", err)
	}

	return writeWord(w, bMatch|bits)
}

func (i B) ObjDump(ctx disasm.ViewCtx) string {
	target := immNum(int64(ctx.Addr()) + i.off.val)
	return "b " + target.textHex()
}
