package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Cbnz — cbnz rt, target.
type Cbnz struct {
	base

	rt     string
	target imm
}

func decodeCbnz(w uint32, addr uint64) Instr {
	return Cbnz{
		base:   newBase(addr, w),
		rt:     armRegName(w&0x1f, w>>31&1 == 1),
		target: immNum(int64(addr) + signExtendN(w>>5&0x7ffff, 19)*4),
	}
}

func (i Cbnz) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("cbnz %s, %s", i.rt, i.target.textHex())
}

func (i Cbnz) Encode(w io.Writer, pc uint64) (int64, error) {
	target := i.target.val

	bits, err := brBits(target, int64(pc), 19)
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

	return writeWord(w, match|bits|num)
}

func (i Cbnz) MarshalJSON() ([]byte, error) {
	return i.marshal("cbnz", i.ObjDump(disasm.DefaultViewCtx()), "Branch",
		map[string]any{"Rt": i.rt, "imm19": i.target.val - int64(i.addr)})
}
