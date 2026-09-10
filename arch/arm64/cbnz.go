package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Cbnz — cbnz rt, off.
type Cbnz struct {
	base

	rt  string
	off imm // pc-relative byte offset
}

// newCbnz - the Cbnz constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newCbnz(b base, rt Reg, off int64) (Cbnz, error) {
	err := requireClass(
		rt,
		"Cbnz",
		"rt",
		"x/w register (register 31 reads as zr — use XZR/WZR)",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return Cbnz{}, err
	}

	return Cbnz{
		base: b,
		rt:   rt.name(),
		off:  immNum(off),
	}, nil
}

func (i Cbnz) ObjDump(ctx disasm.ViewCtx) string {
	target := immNum(int64(ctx.Addr()) + i.off.val)
	return fmt.Sprintf("cbnz %s, %s", i.rt, target.textHex())
}

func (i Cbnz) Encode(w io.Writer) (int64, error) {
	bits, err := offBits(i.off.val, 19)
	if err != nil {
		return 0, fmt.Errorf("cbnz: %w", err)
	}

	num, err := armRegNum(i.rt)
	if err != nil {
		return 0, fmt.Errorf("cbnz: %w", err)
	}

	match := uint32(0x35000000)
	if i.rt[0] == 'x' {
		match = 0xB5000000
	}

	return writeWord(w, match|bits<<5|num)
}

func (Builder) Cbnz(rt Reg, off int64) (Instr, error) {
	return newCbnz(base{}, rt, off)
}

func decodeCbnz(w uint32) (Instr, error) {
	in, err := newCbnz(newBase(w), gprOf(w&0x1f, w>>31&1 == 1), signExtendN(w>>5&0x7ffff, 19)*4)
	if err != nil {
		return nil, err
	}

	return in, nil
}
