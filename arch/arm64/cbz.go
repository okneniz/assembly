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

// Cbz — cbz rt, off: off — the pc-relative byte offset of the branch
// destination (the ±1MB imm19 range is checked at encode time).
// rt — x/w register (register 31 reads as zr — use XZR/WZR).
func (Builder) Cbz(rt Reg, off int64) (Instr, error) {
	if err := requireClass(rt, "Cbz", "rt", "x/w register (register 31 reads as zr — use XZR/WZR)",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	return Cbz{rt: rt.name(), off: immNum(off)}, nil
}

func decodeCbz(w uint32) Instr {
	return Cbz{
		base: newBase(w),
		rt:   armRegName(w&0x1f, w>>31&1 == 1),
		off:  immNum(signExtendN(w>>5&0x7ffff, 19) * 4),
	}
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
