package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Cbz — cbz rt, target (imm19; Rt determines 32/64-bit width: w/x).
type Cbz struct {
	base

	rt  string
	off imm // pc-relative byte offset
}

// newCbz - the Cbz constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newCbz(b base, rt Reg, off int64) (Cbz, error) {
	err := requireClass(
		rt,
		"Cbz",
		"rt",
		"x/w register (register 31 reads as zr — use XZR/WZR)",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return Cbz{}, err
	}

	return Cbz{
		base: b,
		rt:   rt.name(),
		off:  immNum(off),
	}, nil
}

func (i Cbz) ObjDump(ctx disasm.ViewCtx) string {
	target := immNum(int64(ctx.Addr()) + i.off.val)
	return fmt.Sprintf("cbz %s, %s", i.rt, target.textHex())
}

func (i Cbz) Encode(w io.Writer) (int64, error) {
	bits, err := offBits(i.off.val, 19)
	if err != nil {
		return 0, fmt.Errorf("cbz: %w", err)
	}

	num, err := armRegNum(i.rt)
	if err != nil {
		return 0, fmt.Errorf("cbz: %w", err)
	}

	match := uint32(0x34000000) // 32-bit form (w register)
	if i.rt[0] == 'x' {
		match = 0xB4000000
	}

	return writeWord(w, match|bits<<5|num)
}

func (Builder) Cbz(rt Reg, off int64) (Instr, error) {
	return newCbz(base{}, rt, off)
}

func decodeCbz(w uint32) (Instr, error) {
	in, err := newCbz(newBase(w), gprOf(w&0x1f, w>>31&1 == 1), signExtendN(w>>5&0x7ffff, 19)*4)
	if err != nil {
		return nil, err
	}

	return in, nil
}
