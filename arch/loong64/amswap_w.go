package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmswapW - amswap.w rd, rk, rj (3R): rd = old MEM[rj]; MEM[rj] = rk.
type AmswapW struct {
	base

	rd, rk, rj uint8
}

// AmswapW - amswap.w rd, rk, rj.
func (Builder) AmswapW(rd, rk, rj Reg) Instr {
	return AmswapW{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

func decodeAmswapW(w uint32, addr uint64) Instr {
	return AmswapW{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i AmswapW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("amswap.w %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}

func (i AmswapW) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["amswap.w"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i AmswapW) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"amswap.w",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rk": laRegName(i.rk), "rj": laRegName(i.rj)},
	)
}
