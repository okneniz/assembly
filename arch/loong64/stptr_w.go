package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// StptrW - stptr.w rd, rj, offs (DJSk14): MEM[rj + offs] = low32(rd)
// (offs is a word-scaled byte offset, stored raw).
type StptrW struct {
	base

	rd, rj uint8
	off    imm
}

// StptrW - stptr.w rd, rj, offs (the byte offset).
func (Builder) StptrW(rd, rj Reg, off Imm14) Instr {
	return StptrW{
		rd:  rd.Num(),
		rj:  rj.Num(),
		off: immNum(off.Val()),
	}
}

func decodeStptrW(w uint32, addr uint64) Instr {
	return StptrW{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		off:  immNum(sField(w, 10, 14) << 2),
	}
}

func (i StptrW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("stptr.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.off.text())
}

func (i StptrW) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["stptr.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterS(i.off.val>>2, 10, 14)

	return writeWord(w, word)
}

func (i StptrW) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"stptr.w",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "off": i.off.val},
	)
}
