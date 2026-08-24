package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Cbz — cbz rt, target (imm19; Rt determines 32/64-bit width: w/x).
type Cbz struct {
	base

	rt     string
	target imm
}

func decodeCbz(w uint32, addr uint64) Instr {
	return Cbz{
		base:   newBase(addr, w),
		rt:     armRegName(w&0x1f, w>>31&1 == 1),
		target: immNum(int64(addr) + signExtendN(w>>5&0x7ffff, 19)*4),
	}
}

func (i Cbz) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("cbz %s, %s", i.rt, i.target.textHex())
}

func (i Cbz) Encode(w io.Writer, pc uint64) (int64, error) {
	target := i.target.val

	bits, err := brBits(target, int64(pc), 19)
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

	return writeWord(w, match|bits|num)
}

func (i Cbz) MarshalJSON() ([]byte, error) {
	return i.marshal("cbz", i.ObjDump(disasm.DefaultViewCtx()), "Branch",
		map[string]any{"Rt": i.rt, "imm19": i.target.val - int64(i.addr)})
}
