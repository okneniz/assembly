package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmswapDbB - amswap_db.b rd, rk, rj (DJK, @orig_fmt=DKJ - the second assembly
// operand packs the K field [14:10], the third the J field [9:5]):
// rd = AMO(MEM[rj], rk): swap with a built-in barrier, old value to rd.
type AmswapDbB struct {
	base

	rd, rk, rj uint8
}

// NewAmswapDbB - amswap_db.b rd, rk, rj (the operand order the assembler writes).
func NewAmswapDbB(rd, rk, rj Reg) Instr {
	return AmswapDbB{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

func decodeAmswapDbB(w uint32, addr uint64) Instr {
	return AmswapDbB{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i AmswapDbB) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("amswap_db.b %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}

func (i AmswapDbB) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["amswap_db.b"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i AmswapDbB) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"amswap_db.b",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rk": laRegName(i.rk), "rj": laRegName(i.rj)},
	)
}
