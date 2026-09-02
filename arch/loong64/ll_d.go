package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// LlD - ll.d rd, rj, offs (DJSk14): rd = MEM[rj + offs] with a reservation
// (load-linked 64 bits; offs - a byte offset, a multiple of 4, in +-16380).
type LlD struct {
	base

	rd, rj uint8
	off    imm
}

// LlD - ll.d rd, rj, offs (the byte offset from rj).
func (Builder) LlD(rd, rj Reg, off Imm14) Instr {
	return LlD{
		rd:  rd.Num(),
		rj:  rj.Num(),
		off: immNum(off.Val()),
	}
}

func decodeLlD(w uint32, addr uint64) Instr {
	return LlD{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		off:  immNum(sField(w, 10, 14) << 2),
	}
}

func (i LlD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ll.d %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.off.text())
}

func (i LlD) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["ll.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterS(i.off.val>>2, 10, 14)

	return writeWord(w, word)
}

func (i LlD) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"ll.d",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "off": i.off.val},
	)
}
