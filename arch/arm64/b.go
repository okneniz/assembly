package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// B — b target (unconditional branch, imm26, ±128MB).
type B struct {
	base

	target imm // absolute target address
}

const bMatch = 0x14000000

func decodeB(w uint32, addr uint64) Instr {
	return B{
		base:   newBase(addr, w),
		target: immNum(int64(addr) + signExtendN(w&0x3ffffff, 26)*4),
	}
}

func (i B) ObjDump(_ disasm.ViewCtx) string {
	return "b " + i.target.textHex()
}

func (i B) Encode(w io.Writer, pc uint64) (int64, error) {
	target := i.target.val

	bits, err := brBits(target, int64(pc), 26)
	if err != nil {
		return 0, fmt.Errorf("b: %w", err)
	}

	return writeWord(w, bMatch|bits)
}

func (i B) MarshalJSON() ([]byte, error) {
	return i.marshal("b", i.ObjDump(disasm.DefaultViewCtx()), "Branch",
		map[string]any{"imm26": i.target.val - int64(i.addr)})
}
