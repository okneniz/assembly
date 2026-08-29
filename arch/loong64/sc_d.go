package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// ScD - sc.d rd, rj, offs (DJSk14): MEM[rj + offs] = rd conditionally,
// rd = success (store-conditional 64 bits; offs is word-scaled like
// ll.d).
type ScD struct {
	base

	rd, rj uint8
	off    imm
}

// NewScD - sc.d rd, rj, offs (the byte offset from rj).
func NewScD(rd, rj Reg, off Imm14) Instr {
	return ScD{
		rd:  rd.Num(),
		rj:  rj.Num(),
		off: immNum(off.Val()),
	}
}

func decodeScD(w uint32, addr uint64) Instr {
	return ScD{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		off:  immNum(sField(w, 10, 14) << 2),
	}
}

func (i ScD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("sc.d %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.off.text())
}

func (i ScD) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["sc.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterS(i.off.val>>2, 10, 14)

	return writeWord(w, word)
}

func (i ScD) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"sc.d",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "off": i.off.val},
	)
}
