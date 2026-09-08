package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmorDbD - amor_db.d rd, rk, rj (3R): rd = old MEM[rj]; MEM[rj] |= rk. A built-in barrier.
type AmorDbD struct {
	base

	rd, rk, rj uint8
}

// AmorDbD - amor_db.d rd, rk, rj.
func (Builder) AmorDbD(rd, rk, rj Reg) Instr {
	return AmorDbD{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

func decodeAmorDbD(w uint32) Instr {
	return AmorDbD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i AmorDbD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("amor_db.d %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}

func (i AmorDbD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["amor_db.d"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}
