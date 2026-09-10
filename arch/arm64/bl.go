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

// newBl - the Bl constructor: the struct is assembled only
// here (the Builder method and the decoder call it).
func newBl(b base, off imm) (Bl, error) { //nolint:unparam // uniform (Instr, error) decodeCtor type
	return Bl{
		base: b,
		off:  off,
	}, nil
}

const blMatch = 0x94000000

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

// Bl — bl off: off — the pc-relative byte offset of the call
// destination (the ±128MB imm26 range is checked at encode time; the
// absolute target is off + the instruction address).
func (Builder) Bl(off int64) Instr {
	in, err := newBl(base{}, immNum(off))
	if err != nil {
		panic(err) // a plain offset cannot fail
	}

	return in
}

func decodeBl(w uint32) (Instr, error) {
	in, err := newBl(newBase(w), immNum(signExtendN(w&0x3ffffff, 26)*4))
	if err != nil {
		return nil, err
	}

	return in, nil
}
