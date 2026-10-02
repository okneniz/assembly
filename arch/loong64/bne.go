package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Bne - bne rj, rd, offs (2RI16, the manual order swaps the registers):
// if rj != rd, jump to pc + offs (offs is word-scaled). The decoded form
// stores the absolute target.
type Bne struct {
	base

	rd, rj uint8
	off    imm
}

func (i Bne) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("bne %s, %s, %s", laRegName(i.rj), laRegName(i.rd), i.off.text())
}

func (i Bne) Encode(w io.Writer) (int64, error) {
	off, err := encPs2(i.off.val, 16, "bne offset")
	if err != nil {
		return 0, err
	}

	word := loongEncodings["bne"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterS(off, 10, 16)

	return writeWord(w, word)
}
