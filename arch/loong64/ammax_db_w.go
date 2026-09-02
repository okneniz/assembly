package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmmaxDbW - ammax_db.w rd, rk, rj (3R): rd = old MEM[rj]; MEM[rj] = max(MEM[rj], rk), signed. A built-in barrier.
type AmmaxDbW struct {
	base

	rd, rk, rj uint8
}

// AmmaxDbW - ammax_db.w rd, rk, rj.
func (Builder) AmmaxDbW(rd, rk, rj Reg) Instr {
	return AmmaxDbW{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

func decodeAmmaxDbW(w uint32, addr uint64) Instr {
	return AmmaxDbW{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i AmmaxDbW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ammax_db.w %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}

func (i AmmaxDbW) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["ammax_db.w"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i AmmaxDbW) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"ammax_db.w",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rk": laRegName(i.rk), "rj": laRegName(i.rj)},
	)
}
