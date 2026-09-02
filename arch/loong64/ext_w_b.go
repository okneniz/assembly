package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// ExtWB - ext.w.b rd, rj (2R): rd = sign-extend of the low byte of rj to the native width.
type ExtWB struct {
	base

	rd, rj uint8
}

// ExtWB - ext.w.b rd, rj.
func (Builder) ExtWB(rd, rj Reg) Instr {
	return ExtWB{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

func decodeExtWB(w uint32, addr uint64) Instr {
	return ExtWB{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i ExtWB) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ext.w.b %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i ExtWB) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["ext.w.b"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i ExtWB) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"ext.w.b",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj)},
	)
}
