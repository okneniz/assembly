package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// ExtWH - ext.w.h rd, rj (2R): rd = sign-extend of the low half of rj to the native width.
type ExtWH struct {
	base

	rd, rj uint8
}

// NewExtWH - ext.w.h rd, rj.
func NewExtWH(rd, rj Reg) Instr {
	return ExtWH{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

func decodeExtWH(w uint32, addr uint64) Instr {
	return ExtWH{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i ExtWH) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ext.w.h %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i ExtWH) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["ext.w.h"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i ExtWH) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"ext.w.h",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj)},
	)
}
