package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmaddDbB - amadd_db.b rd, rk, rj (3R): rd = old MEM[rj]; MEM[rj] += rk. A built-in barrier.
type AmaddDbB struct {
	base

	rd, rk, rj uint8
}

// AmaddDbB - amadd_db.b rd, rk, rj.
func (Builder) AmaddDbB(rd, rk, rj Reg) Instr {
	return AmaddDbB{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

func decodeAmaddDbB(w uint32, addr uint64) Instr {
	return AmaddDbB{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i AmaddDbB) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("amadd_db.b %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}

func (i AmaddDbB) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["amadd_db.b"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i AmaddDbB) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"amadd_db.b",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rk": laRegName(i.rk), "rj": laRegName(i.rj)},
	)
}
