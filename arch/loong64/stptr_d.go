package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// StptrD - stptr.d rd, rj, offs (DJSk14): MEM[rj + offs] = rd (offs is
// a word-scaled byte offset, stored raw).
type StptrD struct {
	base

	rd, rj uint8
	off    imm
}

// NewStptrD - stptr.d rd, rj, offs (the byte offset).
func NewStptrD(rd, rj Reg, off Imm14) Instr {
	return StptrD{
		rd:  rd.Num(),
		rj:  rj.Num(),
		off: immNum(off.Val()),
	}
}

func decodeStptrD(w uint32, addr uint64) Instr {
	return StptrD{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		off:  immNum(sField(w, 10, 14) << 2),
	}
}

func (i StptrD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("stptr.d %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.off.text())
}

func (i StptrD) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["stptr.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterS(i.off.val>>2, 10, 14)

	return writeWord(w, word)
}

func (i StptrD) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"stptr.d",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "off": i.off.val},
	)
}
