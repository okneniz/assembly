package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// ScW - sc.w rd, rj, offs (DJSk14): MEM[rj + offs] = low32(rd)
// conditionally, rd = success (store-conditional 32 bits; offs is
// word-scaled like ll.w).
type ScW struct {
	base

	rd, rj uint8
	off    imm
}

// NewScW - sc.w rd, rj, offs (the byte offset from rj).
func NewScW(rd, rj Reg, off Imm14) Instr {
	return ScW{
		rd:  rd.Num(),
		rj:  rj.Num(),
		off: immNum(off.Val()),
	}
}

func decodeScW(w uint32, addr uint64) Instr {
	return ScW{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		off:  immNum(sField(w, 10, 14) << 2),
	}
}

func (i ScW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("sc.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.off.text())
}

func (i ScW) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["sc.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterS(i.off.val>>2, 10, 14)

	return writeWord(w, word)
}

func (i ScW) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"sc.w",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "off": i.off.val},
	)
}
