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
	rd, rj, rk uint8
}

func (i AddW) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["add.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i AddW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("add.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}
