package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Bl — bl off (call, imm26).
type Bl struct {
	off imm // pc-relative byte offset
}

// newBl - the Bl constructor: the struct is assembled only
// here (the Builder method and the decoder call it).
func newBl(off imm) (Bl, error) { //nolint:unparam // uniform (Instr, error) decodeCtor type
	return Bl{
		off: off,
	}, nil
}

const blMatch = 0x94000000

func (i Bl) Encode(w io.Writer) (int64, error) {
	bits, err := offBits(i.off.val, 26)
	if err != nil {
		return 0, fmt.Errorf("bl: %w", err)
	}

	return writeWord(w, blMatch|bits)
}

func (i Bl) ObjDump(ctx disasm.ViewCtx) string {
	target := immNum(int64(ctx.Addr()) + i.off.val)
	return "bl " + target.textHex()
}
