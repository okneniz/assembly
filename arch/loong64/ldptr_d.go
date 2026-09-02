package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// LdptrD - ldptr.d rd, rj, offs (DJSk14): rd = MEM[rj + offs] (offs is
// a word-scaled byte offset, stored raw).
type LdptrD struct {
	base

	rd, rj uint8
	off    imm
}

// LdptrD - ldptr.d rd, rj, offs (the byte offset).
func (Builder) LdptrD(rd, rj Reg, off Imm14) Instr {
	return LdptrD{
		rd:  rd.Num(),
		rj:  rj.Num(),
		off: immNum(off.Val()),
	}
}

func decodeLdptrD(w uint32, addr uint64) Instr {
	return LdptrD{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		off:  immNum(sField(w, 10, 14) << 2),
	}
}

func (i LdptrD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ldptr.d %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.off.text())
}

func (i LdptrD) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["ldptr.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterS(i.off.val>>2, 10, 14)

	return writeWord(w, word)
}

func (i LdptrD) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"ldptr.d",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "off": i.off.val},
	)
}
