package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// LdptrW - ldptr.w rd, rj, offs (DJSk14): rd = sign32(MEM[rj + offs])
// (offs is a word-scaled byte offset, stored raw).
type LdptrW struct {
	base

	rd, rj uint8
	off    imm
}

// NewLdptrW - ldptr.w rd, rj, offs (the byte offset).
func NewLdptrW(rd, rj Reg, off Imm14) Instr {
	return LdptrW{
		rd:  rd.Num(),
		rj:  rj.Num(),
		off: immNum(off.Val()),
	}
}

func decodeLdptrW(w uint32, addr uint64) Instr {
	return LdptrW{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		off:  immNum(sField(w, 10, 14) << 2),
	}
}

func (i LdptrW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ldptr.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.off.text())
}

func (i LdptrW) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["ldptr.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterS(i.off.val>>2, 10, 14)

	return writeWord(w, word)
}

func (i LdptrW) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"ldptr.w",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "off": i.off.val},
	)
}
