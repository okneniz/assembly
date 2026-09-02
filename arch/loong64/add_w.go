// Package loong64 — per-instruction LoongArch (LA64) structs: word
// decoders, formatters (objdump notation) and instruction constructors
// (the Builder methods, the exact inverse of decode), with the encoding
// bits joined from the generated loongEncodings table.
package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AddW - add.w rd, rj, rk (3R): rd = sign32(rj + rk).
type AddW struct {
	base

	rd, rj, rk uint8
}

// AddW - add.w rd, rj, rk.
func (Builder) AddW(rd, rj, rk Reg) Instr {
	return AddW{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeAddW(w uint32, addr uint64) Instr {
	return AddW{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i AddW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("add.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i AddW) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["add.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i AddW) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"add.w",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
