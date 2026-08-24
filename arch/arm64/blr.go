package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Blr — blr xn (indirect call).
type Blr struct {
	base

	rn string
}

const blrMatch = 0xD63F0000

func decodeBlr(w uint32, addr uint64) Instr {
	return Blr{
		base: newBase(addr, w),
		rn:   regNameX(w >> 5 & 0x1f),
	}
}

func (i Blr) ObjDump(_ disasm.ViewCtx) string {
	return "blr " + i.rn
}

func (i Blr) Encode(w io.Writer, pc uint64) (int64, error) {
	num, err := armRegNum(i.rn)
	if err != nil {
		return 0, fmt.Errorf("blr: %w", err)
	}

	return writeWord(w, blrMatch|num<<5)
}

func (i Blr) MarshalJSON() ([]byte, error) {
	return i.marshal(
		"blr",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Branch",
		map[string]any{"Rn": i.rn},
	)
}
