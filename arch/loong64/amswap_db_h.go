package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmswapDbH - amswap_db.h rd, rk, rj (3R): rd = old MEM[rj]; MEM[rj] = rk. A built-in barrier.
type AmswapDbH struct {
	base

	rd, rk, rj uint8
}

// AmswapDbH - amswap_db.h rd, rk, rj.
func (Builder) AmswapDbH(rd, rk, rj Reg) Instr {
	return AmswapDbH{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

func decodeAmswapDbH(w uint32, addr uint64) Instr {
	return AmswapDbH{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i AmswapDbH) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("amswap_db.h %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}

func (i AmswapDbH) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["amswap_db.h"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i AmswapDbH) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"amswap_db.h",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rk": laRegName(i.rk), "rj": laRegName(i.rj)},
	)
}
