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

func decodeAmswapW(w uint32) Instr {
	return AmswapW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i AmswapW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("amswap.w %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}

func (i AmswapW) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["amswap.w"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}
