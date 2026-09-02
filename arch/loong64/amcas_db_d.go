package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmcasDbD - amcas_db.d rd, rk, rj (3R): if MEM[rj] == rd then MEM[rj] = rk; rd = old MEM[rj]. A built-in barrier.
type AmcasDbD struct {
	base

	rd, rk, rj uint8
}

// AmcasDbD - amcas_db.d rd, rk, rj.
func (Builder) AmcasDbD(rd, rk, rj Reg) Instr {
	return AmcasDbD{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

func decodeAmcasDbD(w uint32, addr uint64) Instr {
	return AmcasDbD{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i AmcasDbD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("amcas_db.d %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}

func (i AmcasDbD) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["amcas_db.d"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i AmcasDbD) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"amcas_db.d",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rk": laRegName(i.rk), "rj": laRegName(i.rj)},
	)
}
