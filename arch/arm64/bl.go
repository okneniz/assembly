package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Bl — bl target (call, imm26).
type Bl struct {
	base

	target imm
}

const blMatch = 0x94000000

func decodeBl(w uint32, addr uint64) Instr {
	return Bl{
		base:   newBase(addr, w),
		target: immNum(int64(addr) + signExtendN(w&0x3ffffff, 26)*4),
	}
}

func (i Bl) ObjDump(_ disasm.ViewCtx) string {
	return "bl " + i.target.textHex()
}

func (i Bl) Encode(w io.Writer, pc uint64) (int64, error) {
	target := i.target.val

	bits, err := brBits(target, int64(pc), 26)
	if err != nil {
		return 0, fmt.Errorf("bl: %w", err)
	}

	return writeWord(w, blMatch|bits)
}

func (i Bl) MarshalJSON() ([]byte, error) {
	return i.marshal("bl", i.ObjDump(disasm.DefaultViewCtx()), "Branch",
		map[string]any{"imm26": i.target.val - int64(i.addr)})
}
