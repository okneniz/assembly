package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// LdHu - ld.hu rd, rj, si12 (2RI12): rd = the halfword at rj + si12, zero-extended.
type LdHu struct {
	base

	rd, rj uint8
	off    imm
}

// LdHu - ld.hu rd, rj, si12.
func (Builder) LdHu(rd, rj Reg, v Imm12) Instr {
	return LdHu{
		rd:  rd.Num(),
		rj:  rj.Num(),
		off: immNum(v.Val()),
	}
}

func decodeLdHu(w uint32, addr uint64) Instr {
	return LdHu{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		off:  immNum(sField(w, 10, 12)),
	}
}

func (i LdHu) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ld.hu %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.off.text())
}

func (i LdHu) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["ld.hu"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterS(i.off.val, 10, 12)

	return writeWord(w, word)
}

func (i LdHu) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"ld.hu",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "off": i.off.val},
	)
}
