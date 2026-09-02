package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmorDbW - amor_db.w rd, rk, rj (3R): rd = old MEM[rj]; MEM[rj] |= rk. A built-in barrier.
type AmorDbW struct {
	base

	rd, rk, rj uint8
}

// AmorDbW - amor_db.w rd, rk, rj.
func (Builder) AmorDbW(rd, rk, rj Reg) Instr {
	return AmorDbW{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

func decodeAmorDbW(w uint32, addr uint64) Instr {
	return AmorDbW{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i AmorDbW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("amor_db.w %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}

func (i AmorDbW) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["amor_db.w"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i AmorDbW) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"amor_db.w",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rk": laRegName(i.rk), "rj": laRegName(i.rj)},
	)
}
