package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmcasDbW - amcas_db.w rd, rk, rj (3R): if MEM[rj] == rd then MEM[rj] = rk; rd = old MEM[rj]. A built-in barrier.
type AmcasDbW struct {
	base

	rd, rk, rj uint8
}

// AmcasDbW - amcas_db.w rd, rk, rj.
func (Builder) AmcasDbW(rd, rk, rj Reg) Instr {
	return AmcasDbW{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

func decodeAmcasDbW(w uint32, addr uint64) Instr {
	return AmcasDbW{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i AmcasDbW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("amcas_db.w %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}

func (i AmcasDbW) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["amcas_db.w"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i AmcasDbW) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"amcas_db.w",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rk": laRegName(i.rk), "rj": laRegName(i.rj)},
	)
}
