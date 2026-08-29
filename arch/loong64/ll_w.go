package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// LlW - ll.w rd, rj, offs (DJSk14): rd = MEM[rj + offs] with a
// reservation (load-linked 32 bits; offs is word-scaled - a byte offset
// multiple of 4 in +-16380, relative to the rj VALUE, like jirl).
type LlW struct {
	base

	rd, rj uint8
	off    imm
}

// NewLlW - ll.w rd, rj, offs (the byte offset from rj).
func NewLlW(rd, rj Reg, off Imm14) Instr {
	return LlW{
		rd:  rd.Num(),
		rj:  rj.Num(),
		off: immNum(off.Val()),
	}
}

func decodeLlW(w uint32, addr uint64) Instr {
	return LlW{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		off:  immNum(sField(w, 10, 14) << 2),
	}
}

func (i LlW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ll.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.off.text())
}

func (i LlW) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["ll.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterS(i.off.val>>2, 10, 14)

	return writeWord(w, word)
}

func (i LlW) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"ll.w",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "off": i.off.val},
	)
}
