package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// StleB - stle.b rd, rj, rk (DJK): store the low byte of rd with a bounds
// check against rk, trapping outside.
type StleB struct {
	base

	rd, rj, rk uint8
}

// StleB - stle.b rd, rj, rk.
func (Builder) StleB(rd, rj, rk Reg) Instr {
	return StleB{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeStleB(w uint32, addr uint64) Instr {
	return StleB{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i StleB) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("stle.b %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i StleB) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["stle.b"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i StleB) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"stle.b",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
