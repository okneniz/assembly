package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Bfm — bfm rd, rn, #immr, #imms (no aliases).
type Bfm struct {
	base

	rd, rn     string
	immr, imms uint32
	isf        bool
}

const (
	bfmX uint32 = 0xB3400000
	bfmW uint32 = 0x33000000
)

func decodeBfmInstr(w uint32, addr uint64) Instr {
	return Bfm{
		base: newBase(addr, w),
		rd:   armRegName(w&0x1f, w>>31&1 == 1),
		rn:   armRegName(w>>5&0x1f, w>>31&1 == 1),
		immr: w >> 16 & 0x3f,
		imms: w >> 10 & 0x3f,
		isf:  w>>31&1 == 1,
	}
}

func (i Bfm) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("bfm %s, %s, #%d, #%d", i.rd, i.rn, i.immr, i.imms)
}

func (i Bfm) Encode(w io.Writer, pc uint64) (int64, error) {
	return bfmWrite(w, bfmX, bfmW, i.isf, i.rd, i.rn, i.immr, i.imms)
}

func (i Bfm) MarshalJSON() ([]byte, error) {
	return i.marshal("bfm", i.ObjDump(disasm.DefaultViewCtx()), "Data processing - immediate",
		map[string]any{"Rd": i.rd, "Rn": i.rn, "immr": i.immr, "imms": i.imms})
}
