package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Br — br xn (indirect branch).
type Br struct {
	base

	rn string
}

const brMatch = 0xD61F0000

func decodeBr(w uint32, addr uint64) Instr {
	return Br{
		base: newBase(addr, w),
		rn:   regNameX(w >> 5 & 0x1f),
	}
}

func (i Br) ObjDump(_ disasm.ViewCtx) string {
	return "br " + i.rn
}

func (i Br) Encode(w io.Writer, pc uint64) (int64, error) {
	num, err := armRegNum(i.rn)
	if err != nil {
		return 0, fmt.Errorf("br: %w", err)
	}

	return writeWord(w, brMatch|num<<5)
}

func (i Br) MarshalJSON() ([]byte, error) {
	return i.marshal("br", i.ObjDump(disasm.DefaultViewCtx()), "Branch", map[string]any{"Rn": i.rn})
}
