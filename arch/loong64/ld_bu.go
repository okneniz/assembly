package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// LdBu - ld.bu rd, rj, si12 (2RI12): rd = the byte at rj + si12, zero-extended.
type LdBu struct {
	base

	rd, rj uint8
	off    imm
}

// LdBu - ld.bu rd, rj, si12.
func (Builder) LdBu(rd, rj Reg, v Imm12) Instr {
	return LdBu{
		rd:  rd.Num(),
		rj:  rj.Num(),
		off: immNum(v.Val()),
	}
}

func decodeLdBu(w uint32, addr uint64) Instr {
	return LdBu{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		off:  immNum(sField(w, 10, 12)),
	}
}

func (i LdBu) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ld.bu %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.off.text())
}

func (i LdBu) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["ld.bu"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterS(i.off.val, 10, 12)

	return writeWord(w, word)
}

func (i LdBu) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"ld.bu",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "off": i.off.val},
	)
}
