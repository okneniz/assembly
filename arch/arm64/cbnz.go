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

// newCbnz - the Cbnz constructor: the struct is assembled only
// here (the Builder method and the decoder call it).
func newCbnz(b base, rt string, off imm) Cbnz {
	return Cbnz{
		base: b,
		rt:   rt,
		off:  off,
	}
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

// Cbnz — cbnz rt, target: target — the absolute address of the
// branch destination (the ±1MB imm19 range is checked at encode time,
// from pc). rt — x/w register (register 31 reads as zr — use XZR/WZR).
func (Builder) Cbnz(rt Reg, off int64) (Instr, error) {
	if err := requireClass(rt, "Cbnz", "rt", "x/w register (register 31 reads as zr — use XZR/WZR)",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	return newCbnz(base{}, rt.name(), immNum(off)), nil
}

func decodeCbnz(w uint32) Instr {
	return newCbnz(
		newBase(w),
		armRegName(w&0x1f, w>>31&1 == 1),
		immNum(signExtendN(w>>5&0x7ffff, 19)*4),
	)
}
