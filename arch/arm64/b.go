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

// newB - the B constructor: the struct is assembled only
// here (the Builder method and the decoder call it).
func newB(b base, off imm) B {
	return B{
		base: b,
		off:  off,
	}
}

const bMatch = 0x14000000

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

// B — b off (the pc-relative byte offset; the absolute target is off +
// the instruction address).
func (Builder) B(off int64) Instr {
	return newB(base{}, immNum(off))
}

func decodeB(w uint32) Instr {
	return newB(newBase(w), immNum(signExtendN(w&0x3ffffff, 26)*4))
}
